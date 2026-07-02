package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	playground "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator

type RequestError struct {
	err error
}

func (e *RequestError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *RequestError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

func Init() error {
	v, ok := binding.Validator.Engine().(*playground.Validate)
	if !ok {
		return nil
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// {"json", "form", "uri", "header"}
		for _, tagName := range []string{"json", "form"} {
			name := strings.SplitN(fld.Tag.Get(tagName), ",", 2)[0]
			if name != "" && name != "-" {
				return name
			}
		}

		return fld.Name
	})

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	return enTranslations.RegisterDefaultTranslations(v, trans)
}

func IsRequestError(err error) bool {
	var requestErr *RequestError
	return errors.As(err, &requestErr)
}

func NewRequestError(err error) error {
	if err == nil {
		return nil
	}

	return &RequestError{err: err}
}

func formatError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors playground.ValidationErrors
	if errors.As(err, &validationErrors) {
		return NewRequestError(errors.New(formatValidationErrors(validationErrors)))
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		field := unmarshalTypeError.Field
		if field == "" {
			field = "Request Parameter"
		}

		return NewRequestError(fmt.Errorf("%s Type error, should be %s", field, unmarshalTypeError.Type.String()))
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return NewRequestError(errors.New("Request body JSON format error"))
	}

	return err
}

func formatValidationErrors(validationErrors playground.ValidationErrors) string {

	messages := make([]string, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		if trans != nil {
			messages = append(messages, fieldError.Translate(trans))
			continue
		}

		messages = append(messages, fieldError.Error())
	}

	return strings.Join(messages, "；")
}

func Bind(c *gin.Context, obj interface{}) error {
	var err error

	switch c.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
		err = c.ShouldBindQuery(obj)
	default:
		err = c.ShouldBind(obj)
	}

	if err != nil {
		return formatError(err)
	}

	return nil
}
