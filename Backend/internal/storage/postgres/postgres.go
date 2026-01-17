package postgres

import (
	// "fmt"
	"log"
	// "log/slog"

	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/config"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Db *gorm.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connection has been established");
	err = db.AutoMigrate(
		&types.User{},
		&types.Room{},
		&types.RoomMember{},
		&types.Message{},
	)
	if err != nil {
		return nil, err
	}

	return &Postgres{Db: db}, nil
}
func (p *Postgres) RegisterUser(user *types.User) (*types.User, error) {
	if err := p.Db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (p *Postgres) GetUser(user  *types.User)(*types.User,error){
	err:=p.Db.Where(user).First(user).Error;
	if err!=nil{
		return nil,err
	}
	return user,nil
}
func(p *Postgres)RegisterRoom(chatRoom *types.Room)(*types.Room,error){
	if err := p.Db.Create(chatRoom).Error; err != nil {
		return nil, err
	}
	return chatRoom, nil
}
func (p *Postgres) GetRoom(code string) (*types.Room, error) {
    var chatRoom types.Room
    // Specify the column: "room_code = ?"
    if err := p.Db.Where("room_code = ?", code).First(&chatRoom).Error; err != nil {
        return nil, err
    }
    return &chatRoom, nil
}
func (p *Postgres) GetUserRooms(userID uint) ([]types.Room, error) {
    var rooms []types.Room
    err := p.Db.Model(&types.User{ID: userID}).Association("Rooms").Find(&rooms)
    return rooms, err
}
func (p *Postgres) CheckExistingMembers(userid uint, roomid uint) (*types.RoomMember, error) {
    var existingMember types.RoomMember
    err := p.Db.Where("room_id = ? AND user_id = ?", roomid, userid).First(&existingMember).Error
    
    if err != nil {
        // If it's not found, this is GOOD for joining. Return nil.
        return nil, err 
    }

    // If we get here, the user actually exists.
    return &existingMember, nil
}

func (p *Postgres)GetRoomMessages(roomId uint)([]types.Message,error){
	var messages []types.Message
	err := p.Db.Preload("FromUser").
		Where("room_id = ?", roomId).
		Order("created_at asc").
		Find(&messages).Error
		 if err != nil {
        return nil, err 
    }

    return messages, nil
}

func (p *Postgres) GetPrivateMessages(userA, userB uint) ([]types.Message, error) {
    var messages []types.Message

    // We fetch messages where:
    // (A sent to B) OR (B sent to A)
    // AND room_id is null (to ensure it's not a group message)
    err := p.Db.Preload("FromUser").
        Where(
            "(from_id = ? AND to_id = ? AND room_id IS NULL) OR (from_id = ? AND to_id = ? AND room_id IS NULL)",
            userA, userB, userB, userA,
        ).
        Order("created_at asc").
        Find(&messages).Error

    if err != nil {
        return nil, err
    }

    return messages, nil
}