package infrastructure

import (
	"github.com/ThienKim52/golang-dev/internal/api"
	loggerPkg "github.com/ThienKim52/golang-dev/pkg/logger"
	"github.com/ThienKim52/golang-dev/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	utils.NoErr(err)
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg
}

func CreateAPI() api.Engine {
	loggerPkg.SetLogLevel("debug")
	cfg := CreateAPIConfig()
	redisClient := CreateRedisConn()
	sqlDB := CreateSQLWithMigration()
	jwtGen, jwtVal := CreateJWTProvider()
	app := gin.New()
	return api.New(&api.EngineOpts{
		App: app,
		Config: cfg,
		RedisClient: redisClient,
		DB: sqlDB,
		JWTGen: jwtGen,
		JWTVal: jwtVal,
	})
}