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

type Message struct {
    ID        uint      `json:"id" gorm:"primaryKey"` // Needed for DB
    Content   string    `json:"content"`
    
    // Foreign Keys
    RoomID    uint      `json:"room_id"` 
    UserID    uint      `json:"user_id"` 
    
    // Optional: Preload user data (useful for showing sender names in history)
    User      User      `json:"user" gorm:"foreignKey:UserID"`
    
    CreatedAt time.Time `json:"created_at"`
}