package port

import "context"

type NotificationRepository interface {
	SendNotification(ctx context.Context, msg string) error
}
