package mysql

import (
	"fmt"

	"gorm.io/gorm"
)

// MigrateSchema runs migrations for all models
func MigrateSchema(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Tokens{},
		&Users{},
		&Foods{},
		&Categories{},
		&FoodCategories{},
		&Nutrients{},
		// 다른 테이블 추가
	)
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}
	return nil
}
