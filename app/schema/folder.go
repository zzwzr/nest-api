package schema

import (
	"nest-api/app/schema/mixin"
	"nest-api/internal/utils"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Folder struct {
	ent.Schema
}

func (Folder) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "folders"},
	}
}

func (Folder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}

func (Folder) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Comment("主键"),

		field.Int64("project_id").
			Default(0).
			Comment("项目 ID"),

		field.Int64("parent_id").
			Default(0).
			Comment("父文件夹 ID，0 表示根级"),

		field.String("name").
			MaxLen(100).
			Default("").
			Comment("文件夹名称"),

		field.Int("sort_order").
			Default(0).
			Comment("排序"),

		field.Int64("created_by").
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

func (Folder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("folders").
			Unique().
			Required().
			Field("project_id"),
		edge.From("creator", User.Type).
			Ref("created_folders").
			Unique().
			Required().
			Field("created_by"),
		edge.To("interfaces", API.Type),
	}
}
