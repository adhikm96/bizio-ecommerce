package common

import "github.com/Digital-AIR/bizio-ecommerce/internal/model"

type NotificationCreateDto struct {
	UserId           uint   `json:"user_id"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
}

type NotificationListDto struct {
	UserID           uint                     `json:"user_id"`
	ID               uint                     `json:"id"`
	NotificationType string                   `json:"notification_type"`
	Message          string                   `json:"message"`
	Status           model.NotificationStatus `json:"status" `
}
