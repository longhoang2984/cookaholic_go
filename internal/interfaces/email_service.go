package interfaces

import "context"

type EmailService interface {
	SendOTP(ctx context.Context, email, otp string) error
}
