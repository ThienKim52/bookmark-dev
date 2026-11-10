package fixtures

import (
	"testing"

	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

type base struct {
	db *gorm.DB
}

func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

func (b *base) DB() *gorm.DB {
	return b.db
}

func NewFixture(t *testing.T, fix Fixture) Fixture {
	db := sqldb.InitMockDB(t)
	fix.SetupDB(db)
	if err := fix.Migrate(); err != nil {
		t.Fatalf("failed to migrate fixture: %v", err)
	}
	if err := fix.GenerateData(); err != nil {
		t.Fatalf("failed to generate fixture data: %v", err)
	}
	return fix
}