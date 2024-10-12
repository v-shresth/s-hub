package users

import (
	"cms/clients"
	"cms/models"
	"cms/services"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type svc struct {
	log        clients.Logger
	repo       *repo
	tokenMaker clients.TokenMaker
}

func NewUserService(log clients.Logger, systemDb *gorm.DB, tokenMaker clients.TokenMaker) services.UserService {
	return &svc{
		log:        log,
		repo:       newRepo(systemDb, log),
		tokenMaker: tokenMaker,
	}
}

func (s *svc) Register(ctx context.Context, dbUser models.Users) (*models.Users, error) {
	var err error
	dbUser.Password, err = clients.HashPassword(dbUser.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)
	}

	isUserExists, err := s.repo.checkIfUserExists(ctx, dbUser.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot check user existence: %v", err)
	}

	if isUserExists {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	user, err := s.repo.createUser(ctx, dbUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	return &user, nil
}

func (s *svc) FindUserSession(ctx context.Context, sessionId uint) (models.Session, error) {
	return s.repo.findUserSession(ctx, sessionId)
}

func (s *svc) Login(ctx context.Context, email, password string) (*models.Token, error) {
	user, err := s.repo.findUserInfoByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	passwordMatched := clients.CheckPassword(password, user.Password)
	if !passwordMatched {
		return nil, status.Error(codes.FailedPrecondition, "wrong password")
	}

	session, err := s.repo.createUserSession(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	tokenPair, err := s.tokenMaker.GenerateTokenPair(session.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create token pair: %v", err)
	}

	return &tokenPair, nil
}

func (s *svc) Logout(ctx context.Context, sessionId uint) error {
	session, err := s.repo.findUserSession(ctx, sessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Error(codes.NotFound, "session not found")
		}
		return status.Errorf(codes.Internal, "cannot find session: %v", err)
	}

	if !session.EndedAt.Before(time.Now()) {
		return nil
	}

	return s.repo.endUserSession(ctx, sessionId)
}
