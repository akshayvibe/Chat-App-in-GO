package sqlite

import (
	// "database/sql"
	"errors"

	"github.com/akshayjha21/Chat-App-in-GO/internal/config"
	"github.com/akshayjha21/Chat-App-in-GO/internal/types"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	// "golang.org/x/tools/go/analysis/passes/defers"
)

type Sqlite struct {
	Db *gorm.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := gorm.Open(sqlite.Open(cfg.StoragePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&types.User{},
		&types.Room{},
		&types.RoomMember{},
	)
	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}
func (s *Sqlite) CreateConnection(user *types.User, room *types.Room) (*types.RoomMember, error) {

	newMember := types.RoomMember{
		RoomId: room.Id,
		UserId: user.Id,
		Role:   "member",
	}
	result := s.Db.Create(&newMember)
	if result.Error != nil {
		// Specifically check for unique constraint error (user already in room)
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errors.New("user already in this room")
		}
		return nil, result.Error
	}
	return &newMember, nil
}
