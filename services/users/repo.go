package users

import (
	"cms/clients"
	"cms/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type repo struct {
	systemDb *gorm.DB
	log      clients.Logger
}

func newRepo(systemDb *gorm.DB, log clients.Logger) *repo {
	return &repo{
		systemDb: systemDb,
		log:      log,
	}
}

func (r *repo) checkIfUserExists(ctx context.Context, email string) (bool, error) {
	var isUserExists bool
	err := r.systemDb.Debug().WithContext(ctx).Model(&models.Users{}).Where("email = ?", email).Select("count(*)>0").Find(&isUserExists).Error
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("unable to check if user exists")
		return false, err
	}

	return isUserExists, nil
}

func (r *repo) createUser(ctx context.Context, dbUser models.Users) (models.Users, error) {
	err := r.systemDb.Debug().WithContext(ctx).Create(&dbUser).Error
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("unable to create user")
		return dbUser, err
	}

	return dbUser, nil
}

func (r *repo) findUserSession(ctx context.Context, sessionId uint) (models.Session, error) {
	var session models.Session
	err := r.systemDb.Debug().WithContext(ctx).Model(&models.Session{}).
		Where("id=?", sessionId).First(&session).Error
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("unable to check session")
		return session, err
	}

	return session, nil
}

func (r *repo) createUserSession(ctx context.Context, userId uint) (models.Session, error) {
	session := models.Session{
		UserId:    userId,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(24 * time.Hour),
	}
	err := r.systemDb.Debug().WithContext(ctx).Create(&session).Error
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("unable to create session")
		return session, err
	}
	return session, nil
}

func (r *repo) endUserSession(ctx context.Context, sessionId uint) error {
	err := r.systemDb.Debug().WithContext(ctx).Model(&models.Session{}).Where("id=?", sessionId).Update("ended_at", time.Now()).Error
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("unable to end session")
		return err
	}
	return nil
}

func (r *repo) findUserInfoByEmail(ctx context.Context, email string) (models.Users, error) {
	var user models.Users
	err := r.systemDb.Debug().WithContext(ctx).Model(&models.Users{}).
		Where("email=?", email).First(&user).Error
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("unable to find user")
		return user, err
	}

	return user, nil
}
