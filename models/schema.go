package models

import (
	"gorm.io/gorm"
)

type Schema struct {
	SchemaName string
	Fields     []Field
}

type Field struct {
	Name string
	Type string
}

type SchemaMetaData struct {
	gorm.Model
	SystemSchemaName  string `gorm:"column:system_schema_name;type:varchar(250)"`
	DisplaySchemaName string `gorm:"column:display_schema_name;type:varchar(250)"`
	SystemFieldName   string `gorm:"column:system_field_name;type:varchar(250)"`
	DisplayFieldName  string `gorm:"column:display_field_name;type:varchar(250)"`
	DisplayFieldType  string `gorm:"column:display_field_type;type:varchar(250)"`
}

type ListSchemaResponse struct {
	TotalSchemas int
	Schemas      []SchemaDetail
}

type SchemaDetail struct {
	SchemaName   string `gorm:"column:schema_name"`
	NoOfFields   int    `gorm:"column:no_of_fields"`
	TotalSchemas int    `gorm:"column:total_schemas"`
}

type GetSchemaResponse struct {
	MetaData []SchemaMetaData
	Data     []map[string]interface{}
}
