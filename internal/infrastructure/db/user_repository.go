package db

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEntity struct {
	*common.BaseEntity
	Username      string       `json:"-" gorm:"unique;not null"`
	Email         string       `json:"email" gorm:"unique;not null"`
	Password      string       `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	FullName      string       `json:"full_name"`
	EmailVerified bool         `json:"email_verified" gorm:"default:false"`
	OTP           *string      `json:"-" gorm:"default:null"`
	OTPExpiresAt  *time.Time   `json:"-" gorm:"default:null"`
	Avatar        common.Image `json:"avatar" gorm:"serializer:json;type:text;default:null"`
	Bio           string       `json:"bio" gorm:"default:null"`
}

func (UserEntity) TableName() string {
	return "users"
}

func (e *UserEntity) ToDomain() *domain.User {
	return &domain.User{
		BaseModel: &common.BaseModel{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			Status:    e.Status,
		},
		Username:      e.Username,
		Email:         e.Email,
		Password:      e.Password,
		FullName:      e.FullName,
		EmailVerified: e.EmailVerified,
		OTP:           e.OTP,
		OTPExpiresAt:  e.OTPExpiresAt,
		Avatar:        e.Avatar,
		Bio:           e.Bio,
	}
}

// FromDomain converts domain.User to UserEntity
func FromDomain(user *domain.User) *UserEntity {
	return &UserEntity{
		BaseEntity: &common.BaseEntity{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Status:    user.Status,
		},
		Username:      user.Username,
		Email:         user.Email,
		Password:      user.Password,
		FullName:      user.FullName,
		EmailVerified: user.EmailVerified,
		OTP:           user.OTP,
		OTPExpiresAt:  user.OTPExpiresAt,
		Avatar:        user.Avatar,
		Bio:           user.Bio,
	}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(FromDomain(user)).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user UserEntity
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user UserEntity
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user UserEntity
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(FromDomain(user)).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Save(&UserEntity{
		BaseEntity: &common.BaseEntity{
			Status: 0,
		},
	}).Error
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	var users []UserEntity
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}

	userDomains := make([]domain.User, len(users))
	for i := range users {
		userDomains[i] = *users[i].ToDomain()
	}

	return userDomains, nil
}

func (r *userRepository) VerifyOTP(ctx context.Context, id uuid.UUID, otp string) error {
	user, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if user.OTP != nil && *user.OTP != otp {
		return errors.New("OTP not match")
	}

	if user.OTPExpiresAt.Before(time.Now()) {
		return errors.New("OTP has been expired")
	}

	userEntity := FromDomain(user)
	userEntity.OTP = nil
	userEntity.OTPExpiresAt = nil
	userEntity.EmailVerified = true

	return r.db.WithContext(ctx).Where("id = ?", id).Save(userEntity).Error
}