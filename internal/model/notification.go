package model

type Notification struct {
	BaseEntity
	UserID           uint   `json:"user_id" gorm:"index:idx_notifications_user_id; not null"`
	User             User   `gorm:"constraint:OnDelete:CASCADE"`
	NotificationType string `json:"notification_type" gorm:"type:varchar(50)"`
	Message          string `json:"message"`
	Status           int    `json:"status" gorm:"type:int; not null; default:0"`
}
