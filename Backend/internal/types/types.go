package types

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null" valiDate:"required"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Rooms     []Room    `json:"rooms,omitempty" gorm:"many2many:room_members;"`
}
type Room struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" gorm:"unique;not null"`
	RoomCode  string `json:"roomcode" gorm:"unique;not null"`
	Isprivate bool   `json:"is_private" gorm:"default:false"`
	//many users can be in many room
	Members []User `json:"members" gorm:"many2many:room_members;"`
}
type RoomMember struct {
	RoomID uint   `json:"roomid" gorm:"primaryKey"`
	UserID uint   `json:"userid" gorm:"primaryKey"`
	Role   string `json:"role" gorm:"default:'member'"`
}

type Message struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Content   string    `json:"content"`
    
    SenderID  uint      `json:"sender_id"`
    Sender    User      `json:"sender" gorm:"foreignKey:SenderID"`
    
    ReceiverID uint     `json:"receiver_id"`
    Receiver   User     `json:"receiver" gorm:"foreignKey:ReceiverID"`
    
    RoomID    uint      `json:"room_id"`
    CreatedAt time.Time `json:"created_at"`
}