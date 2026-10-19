package fixtures

import (
	"testing"

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

func NewFixture(t *testing. T, fix Fixture) *base {
	return &base{db:fix.DB()}
}