package auth

import (
	"nest-api/configs"
	"nest-api/pkg/auth"
)

func Init() {
	JWT = auth.New(configs.JWT)
}
