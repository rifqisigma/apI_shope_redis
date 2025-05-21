package repository

import (
	"api_shope/dto"
	"api_shope/model"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthRepo interface {
	Register(req *dto.RegisterReq) error
	LoginEmail(email string) (*model.User, error)
}

type authRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAuthRepo(db *gorm.DB, redis *redis.Client) AuthRepo {
	return &authRepo{db, redis}
}

func (r *authRepo) Register(req *dto.RegisterReq) error {
	newUser := model.User{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Name,
	}

	key := fmt.Sprintf("behind:pending:register:%d", newUser.ID)
	message := fmt.Sprintf("Welcome to go journey %v", newUser.Username)
	_, err := r.redis.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, key, map[string]interface{}{
			"id":      newUser.ID,
			"email":   newUser.Email,
			"message": message,
			"op":      "register",
		})
		p.Expire(ctx, key, 10*time.Minute)
		return nil
	})
	if err != nil {
		return err
	}
	return r.db.Model(&model.User{}).Create(&newUser).Error
}

func (r *authRepo) LoginEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Model(&model.User{}).Select("id", "email", "password").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
