package citizen

import (
	"cidadon/internal/region"
	"cidadon/internal/user"

	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	User     user.Model
	UserID   uint
	RegionID uint
	LivesIn  region.Model `gorm:"foreignKey:RegionID"`
}

func (*Model) TableName() string {
	return "citizens"
}
