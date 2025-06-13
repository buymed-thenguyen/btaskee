package cmd

import (
	"btaskee/config"
	"btaskee/db"
	"btaskee/handler"
	"fmt"
	"log"

	dbModel "btaskee/model/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	if cfg == nil {
		panic("empty config")
	}
	fmt.Println("✅ Config loaded")

	dbConn, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	fmt.Println("✅ Connected to DB")

	// Inject DB into domain
	db.InjectDB(dbConn)

	if err = dbModel.AutoMigrateAll(dbConn); err != nil {
		panic(err)
	}
	fmt.Println("✅ Migrated models")

	r := handler.SetupRouter(&cfg.Auth)
	log.Println("🚀 Server running on :", cfg.Port)
	_ = r.Run(":" + cfg.Port)
}
