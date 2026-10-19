package infrastructure

import (
	"github.com/ThienKim52/golang-dev/internal/app/model"
	redisPkg "github.com/ThienKim52/golang-dev/pkg/redis"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/ThienKim52/golang-dev/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateRedisConn() *redis.Client {
	redisClient, err := redisPkg.NewClient("")
	utils.NoErr(err)

	return redisClient
}

func CreateSQLWithMigration() *gorm.DB {
	sqlDB, err := sqldb.NewClient("")
	utils.NoErr(err)
	err = MigrateDB(sqlDB) 
	utils.NoErr(err)
	return sqlDB
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}