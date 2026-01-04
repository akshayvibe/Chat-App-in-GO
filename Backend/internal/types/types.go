package types

import "time"

type User struct {
	ID        uint      `json:"ID" gorm:"primaryKey"`
	Username  string    `json:"Username" gorm:"unique;not null" valiDate:"required"`
	Password  string    `json:"Password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
type Room struct {
	ID        uint   `json:"ID"`
	Name      string `json:"Name" gorm:"unique;not null"`
	Isprivate bool   `json:"Is_private" gorm:"default:false"`
}
type RoomMember struct {
	RoomID uint   `json:"RoomID" gorm:"primaryKey"`
	UserID uint   `json:"UserID" gorm:"primaryKey"`
	Role   string `json:"Role" gorm:"default:'member'"`
}

type Message struct {
	ID      uint   `json:"ID" gorm:"primaryKey"` // Needed for DB
	Content string `json:"content" gorm:"type:text;not null"`

	// Foreign Keys
	RoomID uint `json:"room_ID" gorm:"index"`//indexing helps in db search
	UserID uint `json:"user_ID" gorm:"index"`

	// Optional: Preload user data (useful for showing sender names in history)
	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
}
