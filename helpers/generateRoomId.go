package helpers

import "github.com/google/uuid"

func GenerateRoomId() string {
	return uuid.New().String()
}
