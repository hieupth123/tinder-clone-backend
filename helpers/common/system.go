package common

import (
	"github.com/google/uuid"
)
func GenerateUUID() (s string)  {
	uuidNew, _ := uuid.NewUUID()
	return uuidNew.String()
}
