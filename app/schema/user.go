package schema

import (
	"nest-api/app/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Comment("主键"),

		field.String("name").
			MaxLen(50).
			Default("").
			Comment("用户名"),

		field.String("mobile").
			MaxLen(20).
			Default("").
			Comment("手机号"),

		field.String("password").
			MaxLen(255).
			Default("").
			Sensitive().
			Comment("密码"),

		field.Bool("is_admin").
			Default(false).
			Comment("是否管理员"),

		field.Int8("status").
			Default(1).
			Comment("状态"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}
