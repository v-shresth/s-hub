package utils

import (
	"cms/models"
	"cms/pb"
	"cms/utils/constants"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Converts various input cases (normal case, camel case, etc.) to snake_case
func toSnakeCase(input string) string {
	// Trim leading and trailing spaces
	input = strings.TrimSpace(input)

	// Replace multiple spaces with a single space
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

	// Replace spaces with underscores
	input = strings.ReplaceAll(input, " ", "_")

	// Handle camelCase or PascalCase and convert to snake_case
	var output string
	for i, char := range input {
		if unicode.IsUpper(char) {
			if i > 0 && output[len(output)-1] != '_' {
				output += "_"
			}
			output += string(unicode.ToLower(char))
		} else {
			output += string(char)
		}
	}

	// Ensure no consecutive underscores and lowercase everything
	return strings.ToLower(regexp.MustCompile(`_+`).ReplaceAllString(output, "_"))
}

// Validate name is not a reserved PostgreSQL keyword
func validateReservedKeyword(name string) error {
	// List of reserved PostgreSQL keywords
	reservedKeywords := map[string]struct{}{
		"select": {}, "insert": {}, "update": {}, "delete": {}, "from": {}, "where": {}, "user": {},
		"order": {}, "group": {}, "id": {}, "created_at": {}, "updated_at": {}, "deleted_at": {},
		// Add more reserved keywords here...
	}

	// Check if the name is a reserved keyword
	if _, isReserved := reservedKeywords[strings.ToLower(name)]; isReserved {
		return fmt.Errorf("'%s' is a reserved PostgreSQL keyword", name)
	}
	return nil
}

// Validate that the name is in snake case format
func isValidSnakeCase(name string) bool {
	snakeCaseRegex := regexp.MustCompile(`^[a-z]+(_[a-z]+)*$`)
	return snakeCaseRegex.MatchString(name)
}

// Full validation function that handles multiple input formats and reserved keywords
func validateAndNormalizeName(input string) (string, error) {
	// Convert to snake case
	snakeName := toSnakeCase(input)

	// Check if it's a reserved keyword
	if err := validateReservedKeyword(snakeName); err != nil {
		return "", err
	}

	// Ensure the name is in valid snake case format
	if !isValidSnakeCase(snakeName) {
		return "", fmt.Errorf("'%s' is not in valid snake_case format", snakeName)
	}

	return snakeName, nil
}

func ValidateCreateSchemaRequest(req *pb.CreateSchemaRequest) ([]models.SchemaMetaData, error) {
	var err error

	if req.SchemaName == "" {
		return nil, fmt.Errorf("schema name can't be empty")
	}

	displaySchemaName := req.SchemaName
	req.SchemaName, err = validateAndNormalizeName(req.SchemaName)
	if err != nil {
		return nil, err
	}

	var metaData []models.SchemaMetaData
	for idx, field := range req.Fields {
		if field.Name == "" {
			return nil, fmt.Errorf("field name can't be empty")
		}

		var data models.SchemaMetaData
		data.DisplaySchemaName = displaySchemaName
		data.SystemSchemaName = req.SchemaName
		data.DisplayFieldName = field.Name
		req.Fields[idx].Name, err = validateAndNormalizeName(field.Name)
		if err != nil {
			return nil, err
		}
		data.SystemFieldName = req.Fields[idx].Name
		data.DisplayFieldType = field.Type.String()

		metaData = append(metaData, data)
	}

	return metaData, nil
}

func ConvertCreateSchemaApiReqToDbModel(req *pb.CreateSchemaRequest) (models.Schema, error) {
	schema := models.Schema{
		SchemaName: req.SchemaName,
	}

	var fields []models.Field
	for _, field := range req.Fields {
		f := models.Field{
			Name: field.Name,
			Type: ConvertRequestTypeToSQLType(field.Type, field.Length, field.Precision, field.Scale),
		}

		if f.Type == "" {
			return models.Schema{}, fmt.Errorf("field type is not supported")
		}

		fields = append(fields, f)
	}

	schema.Fields = fields

	return schema, nil
}

func ValidateGetSchemaRequest(req *pb.GetSchemaRequest) (models.Filter, error) {
	if req.SchemaName == "" {
		return models.Filter{}, fmt.Errorf("schema name can't be empty")
	}

	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}

	return models.Filter{
		PageSize:   int(req.PageSize),
		PageNumber: int(req.PageNumber),
	}, nil

}

func ConvertGetSchemaDbRespToApiResp(resp models.GetSchemaResponse) (*pb.GetSchemaResponse, error) {
	var out = &pb.GetSchemaResponse{}
	var err error

	out.Records, err = ConvertDbRecordsToApiRecords(resp.Data, resp.MetaData)
	if err != nil {
		return nil, err
	}

	for _, data := range resp.MetaData {
		out.Fields = append(out.Fields, &pb.Field{
			Name: data.DisplayFieldName,
			Type: pb.Type(pb.Type_value[data.DisplayFieldType]),
		})
	}

	return out, nil
}

func ValidateDropSchemaRequest(req *pb.DropSchemaRequest) error {
	if req.SchemaName == "" {
		return fmt.Errorf("schema name can't be empty")
	}

	if req.SchemaName == constants.MetadataSchema {
		return fmt.Errorf("default tables can't be deleted")
	}

	return nil
}
