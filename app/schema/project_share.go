package schema

import (
	"nest-api/app/schema/mixin"
	"nest-api/internal/utils"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProjectShare struct {
	ent.Schema
}

func (ProjectShare) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}

func (ProjectShare) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Comment("主键"),

		field.Int64("project_id").
			Default(0).
			Comment("项目 ID"),

		field.Int64("workspace_id").
			Default(0).
			Comment("工作空间 ID"),

		field.String("name").
			MaxLen(150).
			Default("").
			Comment("分享名称"),

		field.String("share_code").
			MaxLen(16).
			Unique().
			Comment("分享码"),

		field.Bool("enabled").
			Default(true).
			Comment("是否开启分享"),

		field.String("password").
			MaxLen(255).
			Default("").
			Sensitive().
			Comment("访问密码（哈希，空表示无密码）"),

		field.Uint8("permission").
			Default(1).
			Comment("分享权限：1=仅查看"),

		field.Int64("created_by").
			Default(0).
			Comment("创建者用户 ID"),

		field.Time("created_at").
			GoType(utils.DateTime{}).
			Default(utils.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamp(0) without time zone",
			}).
			Comment("创建时间"),

		field.Time("updated_at").
			GoType(utils.DateTime{}).
			Default(utils.Now).
			UpdateDefault(utils.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamp(0) without time zone",
			}).
			Comment("更新时间"),

		field.Time("deleted_at").
			Optional().
			Nillable().
			GoType(utils.DateTime{}).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamp(0) without time zone",
			}).
			Comment("删除时间"),
	}
}

func (ProjectShare) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("shares").
			Unique().
			Required().
			Field("project_id"),
		edge.From("creator", User.Type).
			Ref("created_project_shares").
			Unique().
			Required().
			Field("created_by"),
		edge.To("items", ProjectShareInterface.Type),
	}
}
