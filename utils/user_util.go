package utils

import (
	"cms/models"
	"cms/pb"
	"errors"
	"strings"
)

func ValidateRegisterRequest(req *pb.RegisterRequest) error {
	if strings.TrimSpace(req.User.FirstName) == "" {
		return errors.New("first name can't be empty")
	}

	if strings.TrimSpace(req.User.LastName) == "" {
		return errors.New("last name can't be empty")
	}

	if strings.TrimSpace(req.User.Email) == "" {
		return errors.New("email can't be empty")
	}

	if strings.TrimSpace(req.User.Password) == "" {
		return errors.New("password can't be empty")
	}

	return nil
}

func ConvertApiRegisterRequestToDbRegisterRequest(req *pb.RegisterRequest) models.Users {
	return models.Users{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Email:     req.User.Email,
		Password:  req.User.Password,
	}
}

func ConvertDbRegisterRequestToApiRegisterRequest(dbUser *models.Users) *pb.User {
	if dbUser == nil {
		return nil
	}
	return &pb.User{
		Id:        int32(dbUser.ID),
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
	}
}

func ValidateLoginRequest(req *pb.LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("first name can't be empty")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password can't be empty")
	}

	return nil
}
