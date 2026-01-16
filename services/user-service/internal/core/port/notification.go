package port

import "context"

type NotificationRepository interface {
	SendEmailNotification(ctx context.Context, email string) error
}
