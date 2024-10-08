package constants

import "cms/pb"

const (
	DefaultPort    = ":8080"
	MetadataSchema = "schema_meta_data"
)

type Mode string

const (
	Local Mode = "local"
	Dev   Mode = "dev"
)

var (
	DefaultMetaDataColumns = []struct {
		SystemFieldName  string
		DisplayFieldName string
		DisplayFieldType string
	}{
		{
			SystemFieldName:  "id",
			DisplayFieldType: pb.Type_TYPE_INTEGER.String(),
			DisplayFieldName: "Id",
		},
		{
			SystemFieldName:  "created_at",
			DisplayFieldType: pb.Type_TYPE_TIMESTAMP.String(),
			DisplayFieldName: "Created At",
		},
		{
			SystemFieldName:  "updated_at",
			DisplayFieldType: pb.Type_TYPE_TIMESTAMP.String(),
			DisplayFieldName: "Updated At",
		},
	}

	SystemDefaultColumnsMap = map[string]struct{}{
		"id":         {},
		"created_at": {},
		"updated_at": {},
		"deleted_at": {},
	}
)
