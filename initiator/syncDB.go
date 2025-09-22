package initiator

import (
	"log"

	models "github.com/TNAHOM/ATS-system-main/internal/constants/model"
	"gorm.io/gorm"
)

func SyncDatabase(db *gorm.DB) {
	// db.Migrator().DropTable(&models.JobPost{})
	var err error
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("User AutoMigrate failed: %v", err)
	}
	err = db.AutoMigrate(&models.JobPost{})
	if err != nil {
		log.Fatalf("JobPost AutoMigrate failed: %v", err)
	}
}
