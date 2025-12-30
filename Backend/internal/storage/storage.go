package storage

import "github.com/akshayjha21/Chat-App-in-GO/internal/types"

// import "os/user"
type Storage interface {
	CreateConnection(user *types.User,room *types.Room) (*types.RoomMember,error)
	RegisterUser(user *types.User)(*types.User,error)
}