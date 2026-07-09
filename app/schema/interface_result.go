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

type InterfaceResult struct {
	ent.Schema
}

func (InterfaceResult) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "interface_results"},
	}
}

func (InterfaceResult) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}

func (InterfaceResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive().Unique().Comment("主键"),
		field.Int64("interface_id").Default(0).Comment("接口 ID"),
		field.String("name").MaxLen(100).Default("").Comment("名称"),
		field.Int("status_code").Default(200).Comment("状态码"),
		field.String("format").MaxLen(20).Default("JSON").Comment("格式"),
		field.String("data_type").MaxLen(20).Default("Object").Comment("数据类型"),
		field.Int("sort_order").Default(0).Comment("排序"),
		field.Time("created_at").GoType(utils.DateTime{}).Default(utils.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamp(0) without time zone"}).Comment("创建时间"),
		field.Time("updated_at").GoType(utils.DateTime{}).Default(utils.Now).UpdateDefault(utils.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamp(0) without time zone"}).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().GoType(utils.DateTime{}).
			SchemaType(map[string]string{dialect.Postgres: "timestamp(0) without time zone"}).Comment("删除时间"),
	}
}

func (InterfaceResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("interface", API.Type).Ref("response_results").Unique().Required().Field("interface_id"),
		edge.To("fields", InterfaceField.Type),
	}
}
