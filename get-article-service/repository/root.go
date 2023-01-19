package repository

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	// ErrNotFound error when record not found
	ErrNotFound = fmt.Errorf("record Not Found")
)

type mysqlDBRepository struct {
	mysql *gorm.DB
}
