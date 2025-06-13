package db

import "gorm.io/gorm"

func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Quiz{},
		&Question{},
		&Participant{},
		&ParticipantAnswer{},
		&AnswerOption{},
		&Session{},
	)
}
