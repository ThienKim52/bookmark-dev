package sqldb
import (
	"testing"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
func InitMockDB(t *testing.T) *gorm.DB {
	dsn := ":memory:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create test db: %v", err)
	}
	
	return db
}