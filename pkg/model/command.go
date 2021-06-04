package model

import (
	"time"

	"gorm.io/gorm"
)

type Command struct {
	ID              int            `json:"id" gorm:"primarykey"`
	CommandContent  string         `json:"command_content" gorm:"type:text"`
	Name            string         `json:"name" gorm:"type:varchar(256)"`
	Type            string         `json:"type" gorm:"type:varchar(30)"`
	Description     string         `json:"description" gorm:"type:varchar(1024)"`
	WorkingDir      string         `json:"working_dir" gorm:"type:varchar(1024)"`
	Timeout         int            `json:"timeout" gorm:"type:int"`
	EnableParameter bool           `json:"enableParameter" gorm:"type:int"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func CreateCommand(cmd *Command) error {
	return db.Save(cmd).Error
}
