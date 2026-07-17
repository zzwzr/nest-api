package schema

import (
	"nest-api/internal/utils"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ProjectShareInterface struct {
	ent.Schema
}

func (ProjectShareInterface) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			Comment("主键"),

		field.Int64("share_id").
			Default(0).
			Comment("分享 ID"),

		field.Int64("interface_id").
			Default(0).
			Comment("接口 ID"),

		field.Time("created_at").
			GoType(utils.DateTime{}).
			Default(utils.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamp(0) without time zone",
			}).
			Comment("创建时间"),
	}
}

func (ProjectShareInterface) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("share", ProjectShare.Type).
			Ref("items").
			Unique().
			Required().
			Field("share_id"),
		edge.From("interface", API.Type).
			Ref("share_items").
			Unique().
			Required().
			Field("interface_id"),
	}
}

func (ProjectShareInterface) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("share_id", "interface_id").Unique(),
	}
}
