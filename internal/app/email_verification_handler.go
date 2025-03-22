package app

import (
	"context"
	"cookaholic/internal/interfaces"
	"fmt"
	"math/rand"
	"time"
)

type EmailVerificationHandler struct {
	userRepo     interfaces.UserRepository
	emailService interfaces.EmailService
}

func NewEmailVerificationHandler(userRepo interfaces.UserRepository, emailService interfaces.EmailService) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (h *EmailVerificationHandler) Handle(ctx context.Context, event interfaces.Event) error {
	userEvent, ok := event.(interfaces.UserCreatedEvent)
	if !ok {
		return nil
	}

	user, err := h.userRepo.FindByID(ctx, userEvent.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return interfaces.ErrUserNotFound
	}

	// Generate 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(5 * time.Minute)

	// Update user with OTP
	user.OTP = &otp
	user.OTPExpiresAt = &expiresAt
	if err := h.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Send OTP via email
	return h.emailService.SendOTP(ctx, user.Email, otp)
}
