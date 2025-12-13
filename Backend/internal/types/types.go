package types

import "time"

type User struct {
	Id        uint      `json:"Id" gorm:"primaryKey"`
	Name      string    `json:"Name" validate:"required"`
	Email     string    `json:"Email" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
}
type Room struct {
	Id        uint   `json:"Id"`
	Name      string `json:"Name"`
	Isprivate bool   `json:"Is_private" gorm:"default:false"`
}
type RoomMember struct {
	RoomId uint   `json:"RoomID" gorm:"primaryKey"`
	UserId uint   `json:"UserID" gorm:"primaryKey"`
	Role   string `json:"Role" gorm:"default:'member'"`
}
