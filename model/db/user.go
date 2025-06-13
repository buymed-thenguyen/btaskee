package db

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"` // sẽ lưu password đã mã hóa (hash)
}
