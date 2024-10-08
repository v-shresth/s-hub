package utils

import (
	"cms/pb"
	"cms/utils/constants"
	"errors"
	"fmt"
)

func ValidateCreateRecordRequest(req *pb.CreateRecordRequest) error {
	if req.SchemaName == "" {
		return errors.New("schema name is required")
	}

	for _, record := range req.Records {
		for field, _ := range record.Values {
			if field == "" {
				return errors.New("record fields can't be empty")
			}

			normalizeField, err := validateAndNormalizeName(field)
			if err != nil {
				return err
			}

			if _, ok := constants.SystemDefaultColumnsMap[normalizeField]; ok {
				return fmt.Errorf("you cannot set value to default parameters %s", normalizeField)
			}
		}
	}

	return nil
}

func ValidateGetRecordRequest(req *pb.GetRecordRequest) error {
	if req.SchemaName == "" {
		return errors.New("schema name is required")
	}

	if req.RecordId == 0 {
		return errors.New("record id is required")
	}

	return nil
}

func ValidateDeleteRecordRequest(req *pb.DeleteRecordRequest) error {
	if req.SchemaName == "" {
		return errors.New("schema name is required")
	}

	if req.RecordId == 0 {
		return errors.New("record id is required")
	}

	return nil
}

func ValidateUpdateRecordRequest(req *pb.UpdateRecordRequest) error {
	if req.SchemaName == "" {
		return errors.New("schema name is required")
	}

	if req.RecordId == 0 {
		return errors.New("record id is required")
	}

	for field, _ := range req.Record.Values {
		if field == "" {
			return errors.New("record fields can't be empty")
		}

		normalizeField, err := validateAndNormalizeName(field)
		if err != nil {
			return err
		}

		if _, ok := constants.SystemDefaultColumnsMap[normalizeField]; ok {
			return fmt.Errorf("you cannot update value to default parameters %s", normalizeField)
		}
	}

	return nil
}
