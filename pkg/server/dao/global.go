package dao

import (
	"gorm.io/gorm"
	"vocab-builder/pkg/server/conf"
)

func GetDB() *gorm.DB {
	return conf.DB.Session(&gorm.Session{})
}
