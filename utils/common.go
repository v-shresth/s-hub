package utils

import (
	"cms/models"
	"cms/pb"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ConvertRequestTypeToSQLType(fieldType pb.Type, length int32, precision int32, scale int32) string {
	switch fieldType {
	case pb.Type_TYPE_TEXT:
		return "TEXT"
	case pb.Type_TYPE_VARCHAR:
		// Default length for VARCHAR if not provided
		if length <= 0 {
			length = 255 // Default VARCHAR length
		}
		return fmt.Sprintf("VARCHAR(%d)", length)
	case pb.Type_TYPE_CHAR:
		if length <= 0 {
			length = 1 // Default CHAR length
		}
		return fmt.Sprintf("CHAR(%d)", length)
	case pb.Type_TYPE_INTEGER:
		return "INTEGER"
	case pb.Type_TYPE_BIGINT:
		return "BIGINT"
	case pb.Type_TYPE_SMALLINT:
		return "SMALLINT"
	case pb.Type_TYPE_BOOLEAN:
		return "BOOLEAN"
	case pb.Type_TYPE_DATE:
		return "DATE"
	case pb.Type_TYPE_TIMESTAMP:
		return "TIMESTAMP"
	case pb.Type_TYPE_NUMERIC:
		if precision <= 0 {
			precision = 10 // Default precision
		}
		if scale <= 0 {
			scale = 2 // Default scale
		}
		return fmt.Sprintf("NUMERIC(%d, %d)", precision, scale)
	default:
		return ""
	}
}

func ConvertDbRecordsToApiRecords(dbRecords []map[string]interface{}, metaData []models.SchemaMetaData) ([]*pb.Record, error) {
	// Create a map for faster lookup by frontend name
	metaMap := make(map[string]models.SchemaMetaData)
	for _, data := range metaData {
		metaMap[data.SystemFieldName] = data
	}

	var apiRecords []*pb.Record
	// Iterate over each record in the request
	for _, record := range dbRecords {
		apiRecord := &pb.Record{
			Values: make(map[string]*pb.Value),
		}
		// Iterate over each field in the record
		for columnName, value := range record {
			// Lookup metadata for the frontend column name
			meta, ok := metaMap[columnName]
			if !ok {
				continue // Skip if no metadata found
			}

			columnName = meta.DisplayFieldName // Use frontend field name for API response

			switch meta.DisplayFieldType {
			case pb.Type_TYPE_TEXT.String():
				if v, ok := value.(string); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_TextValue{
							TextValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type TEXT", value))
				}

			case pb.Type_TYPE_VARCHAR.String():
				if v, ok := value.(string); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_VarcharValue{
							VarcharValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type VARCHAR", value))
				}

			case pb.Type_TYPE_CHAR.String():
				if v, ok := value.(string); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_CharValue{
							CharValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type CHAR", value))
				}

			case pb.Type_TYPE_INTEGER.String():
				if v, ok := value.(int32); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_IntValue{
							IntValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type INTEGER", value))
				}

			case pb.Type_TYPE_SMALLINT.String():
				if v, ok := value.(int32); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_SmallintValue{
							SmallintValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type SMALLINT", value))
				}

			case pb.Type_TYPE_BIGINT.String():
				if v, ok := value.(int64); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_BigintValue{
							BigintValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type BIGINT", value))
				}

			case pb.Type_TYPE_BOOLEAN.String():
				if v, ok := value.(bool); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_BoolValue{
							BoolValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type BOOLEAN", value))
				}

			case pb.Type_TYPE_DATE.String():
				if v, ok := value.(time.Time); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_DateValue{
							DateValue: timestamppb.New(v),
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type DATE", value))
				}

			case pb.Type_TYPE_TIMESTAMP.String():
				if v, ok := value.(time.Time); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_TimestampValue{
							TimestampValue: timestamppb.New(v),
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type TIMESTAMP", value))
				}

			case pb.Type_TYPE_NUMERIC.String():
				if v, ok := value.(float64); ok || value == nil {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_NumericValue{
							NumericValue: float32(v),
						},
					}
				} else if v, ok := value.(float32); ok {
					apiRecord.Values[columnName] = &pb.Value{
						Value: &pb.Value_NumericValue{
							NumericValue: v,
						},
					}
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type NUMERIC", value))
				}

			default:
				return nil, errors.New(fmt.Sprintf("Unsupported field type: %s", meta.DisplayFieldType))
			}
		}
		// Append the validated record to the final result
		apiRecords = append(apiRecords, apiRecord)
	}

	return apiRecords, nil
}

func ConvertApiRecordsToDbRecords(apiRecords []*pb.Record, metaData []models.SchemaMetaData) ([]map[string]interface{}, error) {
	// Create a map for faster lookup by backend name
	metaMap := make(map[string]models.SchemaMetaData)
	for _, data := range metaData {
		metaMap[data.DisplayFieldName] = data
	}

	var dbRecords []map[string]interface{}

	// Iterate over each API record
	for _, record := range apiRecords {
		dbRecord := make(map[string]interface{})
		// Iterate over each field in the API record
		for columnName, value := range record.Values {
			// Lookup metadata for the backend column name
			meta, ok := metaMap[columnName]
			if !ok {
				return nil, errors.New(fmt.Sprintf("column %s not found", columnName))
			}

			systemFieldName := meta.SystemFieldName // Use system field name for DB insertion

			switch v := value.Value.(type) {
			case *pb.Value_TextValue:
				if meta.DisplayFieldType == pb.Type_TYPE_TEXT.String() {
					dbRecord[systemFieldName] = v.TextValue
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type %s", v.TextValue, meta.DisplayFieldType))
				}

			case *pb.Value_VarcharValue:
				if meta.DisplayFieldType == pb.Type_TYPE_VARCHAR.String() {
					dbRecord[systemFieldName] = v.VarcharValue
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type %s", v.VarcharValue, meta.DisplayFieldType))
				}

			case *pb.Value_CharValue:
				if meta.DisplayFieldType == pb.Type_TYPE_CHAR.String() {
					dbRecord[systemFieldName] = v.CharValue
				} else {
					return nil, errors.New(fmt.Sprintf("%s is not of type CHAR", v.CharValue))
				}

			case *pb.Value_IntValue:
				if meta.DisplayFieldType == pb.Type_TYPE_INTEGER.String() {
					dbRecord[systemFieldName] = v.IntValue
				} else {
					return nil, errors.New(fmt.Sprintf("%d is not of type INTEGER", v.IntValue))
				}

			case *pb.Value_SmallintValue:
				if meta.DisplayFieldType == pb.Type_TYPE_SMALLINT.String() {
					dbRecord[systemFieldName] = v.SmallintValue
				} else {
					return nil, errors.New(fmt.Sprintf("%d is not of type SMALLINT", v.SmallintValue))
				}

			case *pb.Value_BigintValue:
				if meta.DisplayFieldType == pb.Type_TYPE_BIGINT.String() {
					dbRecord[systemFieldName] = v.BigintValue
				} else {
					return nil, errors.New(fmt.Sprintf("%d is not of type BIGINT", v.BigintValue))
				}

			case *pb.Value_BoolValue:
				if meta.DisplayFieldType == pb.Type_TYPE_BOOLEAN.String() {
					dbRecord[systemFieldName] = v.BoolValue
				} else {
					return nil, errors.New(fmt.Sprintf("%t is not of type BOOLEAN", v.BoolValue))
				}

			case *pb.Value_DateValue:
				if meta.DisplayFieldType == pb.Type_TYPE_DATE.String() {
					dbRecord[systemFieldName] = v.DateValue.AsTime()
				} else {
					return nil, errors.New(fmt.Sprintf("%v is not of type DATE", v.DateValue.AsTime()))
				}

			case *pb.Value_TimestampValue:
				if meta.DisplayFieldType == pb.Type_TYPE_DATE.String() {
					dbRecord[systemFieldName] = v.TimestampValue.AsTime()
				} else {
					return nil, errors.New(fmt.Sprintf("%v is not of type DATE", v.TimestampValue.AsTime()))
				}

			case *pb.Value_NumericValue:
				if meta.DisplayFieldType == pb.Type_TYPE_NUMERIC.String() {
					dbRecord[systemFieldName] = v.NumericValue
				} else {
					return nil, errors.New(fmt.Sprintf("%f is not of type NUMERIC", v.NumericValue))
				}

			default:
				return nil, errors.New(fmt.Sprintf("Unsupported field type for column: %s", systemFieldName))
			}
		}

		// Append the validated dbRecord to the final result
		dbRecords = append(dbRecords, dbRecord)
	}

	return dbRecords, nil
}
