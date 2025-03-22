package app

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo     interfaces.UserRepository
	eventBus interfaces.EventBus
}

func NewUserService(repo interfaces.UserRepository, eventBus interfaces.EventBus) *UserService {
	return &UserService{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *UserService) Create(ctx context.Context, input interfaces.CreateUserInput) (*domain.User, error) {
	// Check if email exists
	if existing, _ := s.repo.FindByEmail(ctx, input.Email); existing != nil {
		return nil, interfaces.ErrEmailExists
	}

	// Check if username exists
	if existing, _ := s.repo.FindByUsername(ctx, input.Username); existing != nil {
		return nil, interfaces.ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		BaseModel: &common.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		},
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Publish user created event
	event := interfaces.UserCreatedEvent{
		UserID: user.ID,
		Email:  user.Email,
	}
	if err := s.eventBus.Publish(ctx, event); err != nil {
		log.Printf("Failed to publish user created event: %v", err)
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, interfaces.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, input interfaces.UpdateUserInput) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, interfaces.ErrUserNotFound
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if input.Avatar != nil {
		user.Avatar = *input.Avatar
	}

	if input.Bio != "" {
		user.Bio = input.Bio
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return interfaces.ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context, page, pageSize int) ([]domain.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}

func (s *UserService) ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, interfaces.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, interfaces.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) VerifyOTP(ctx context.Context, id uuid.UUID, otp string) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return interfaces.ErrUserNotFound
	}

	return s.repo.VerifyOTP(ctx, id, otp)
}

func (s *UserService) ResendOTP(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return interfaces.ErrUserNotFound
	}

	// Publish user created event
	event := interfaces.UserCreatedEvent{
		UserID: user.ID,
		Email:  user.Email,
	}
	if err := s.eventBus.Publish(ctx, event); err != nil {
		log.Printf("Failed to publish user created event: %v", err)
		return err
	}

	return nil
}
