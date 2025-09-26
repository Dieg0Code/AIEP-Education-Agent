package models

import "gorm.io/gorm"

// AutoMigrateAll ejecuta las migraciones de todos los modelos
func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&ChatSession{},
		&ChatMessage{},
		&Insight{},
		&Module{},
		&Topic{},
		&Enrollment{},
	)
}
