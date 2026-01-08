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
func (p *Postgres)CheckExistingMembers(userid uint,roomid uint)(*types.RoomMember,error){
	var existingMember types.RoomMember
    err := p.Db.Where("room_id = ? AND user_id = ?", roomid, userid).First(&existingMember).Error
    
    if err == nil {
        // No error means a record was found
        return nil,err
    }
	return &existingMember,nil
}

