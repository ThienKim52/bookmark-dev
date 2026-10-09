package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/api/middleware"
	healthcheckHdl "github.com/ThienKim52/golang-dev/internal/app/handler/healthcheck"
	userSvc "github.com/ThienKim52/golang-dev/internal/app/service/user"
	userHdl "github.com/ThienKim52/golang-dev/internal/app/handler/user"
	repoLink "github.com/ThienKim52/golang-dev/internal/app/repository/link"
	repo "github.com/ThienKim52/golang-dev/internal/app/repository/user"
	genpassHdl"github.com/ThienKim52/golang-dev/internal/app/handler/genpass"
	genpassSvc"github.com/ThienKim52/golang-dev/internal/app/service/genpass"
	linkHdl "github.com/ThienKim52/golang-dev/internal/app/handler/link"
	linkSvc "github.com/ThienKim52/golang-dev/internal/app/service/link"
	healthcheckSvc "github.com/ThienKim52/golang-dev/internal/app/service/healthcheck"
	"github.com/ThienKim52/golang-dev/pkg/jwtutils"
	"github.com/ThienKim52/golang-dev/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var swaggerOnce sync.Once

// interface Engine
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app         *gin.Engine
	config      *Config
	redisClient *redis.Client
	db          *gorm.DB
	jwtGen      jwtutils.JWTGenerator
	jwtVal      jwtutils.JWTValidator
}

// Constructor
func NewEngine(config *Config, redisClient *redis.Client) Engine {
	app := &engine{
		app:         gin.Default(),
		config:      config,
		redisClient: redisClient,
	}
	app.initRoutes()
	return app
}

type EngineOpts struct {
	App         *gin.Engine
	Config      *Config
	RedisClient *redis.Client
	DB          *gorm.DB
	JWTGen      jwtutils.JWTGenerator
	JWTVal      jwtutils.JWTValidator
}


func New(opts *EngineOpts) Engine {
	app := &engine{
		app:         opts.App,
		config:      opts.Config,
		redisClient: opts.RedisClient,
		db: opts.DB,
		jwtGen: opts.JWTGen,
		jwtVal: opts.JWTVal,
	}
	app.initRoutes()
	return app
}

// Start server
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.config.AppPort))
}

// ServeHTTP implements http.Handler interface
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initHandlers() *allHandlers{
	linkRepo := repoLink.NewRedisLinkRepository(e.redisClient)

	// Health check endpoint with Redis check
	healthCheckSvc := healthcheckSvc.NewHealthCheck(linkRepo, e.config.ServiceName, e.config.InstanceID)
	healthCheckSvcHdl := healthcheckHdl.NewHealthCheck(healthCheckSvc)
	// Genpass endpoint
	genPassSvc := genpassSvc.NewGenPass()
	genPassHdl := genpassHdl.NewGenPass(genPassSvc)
	// URL shortening endpoint
	linkSvc := linkSvc.NewLinkService(linkRepo, genPassSvc)
	linkHdl := linkHdl.NewLinkHandler(linkSvc)
	// Register user endpoint
	userRepo := repo.NewSqlRepository(e.db)
	hasher := utils.NewHasher()
	userService := userSvc.NewService(userRepo, hasher, e.jwtGen)
	userHdl := userHdl.NewHandler(userService)

	return &allHandlers{
		healthCheckSvcHdl: healthCheckSvcHdl,
		genPassHdl:        genPassHdl,
		linkHdl:           linkHdl,
		userHdl:           userHdl,
	}
}

type allHandlers struct {
	healthCheckSvcHdl healthcheckHdl.HealthCheckHandler
	genPassHdl        genpassHdl.GenPass
	linkHdl           linkHdl.LinkHandler
	userHdl           userHdl.Handler
}

// Initialize routes
func (e *engine) initRoutes() {
	// Initialize repository with Redis client
	allHandlers := e.initHandlers()
	jwtAuth := middleware.NewJWTAuth(e.jwtVal)
	swaggerOnce.Do(func() {
		docs.SwaggerInfo.BasePath = e.config.BasePath
	})
	e.app.GET("/health-check", allHandlers.healthCheckSvcHdl.GetResponse)

	e.app.GET("/genpass", allHandlers.genPassHdl.GeneratePassword)

	// e.app.POST("/test/:id/:name", allHandlers.userHdl.Test)
	v1Routes := e.app.Group("/v1")
	{
		// Link-related
		v1Routes.POST("/links/shorten", allHandlers.linkHdl.ShortenURL)
		v1Routes.GET("/links/redirect/:code", allHandlers.linkHdl.Redirect)
		v1Routes.POST("/users/register", allHandlers.userHdl.Register)
		v1Routes.POST("/users/login", allHandlers.userHdl.Login)
		
	}
	selfRoutes := v1Routes.Group("/self")
	selfRoutes.Use(jwtAuth.JWTAuth())
	{
		selfRoutes.GET("/info", allHandlers.userHdl.GetUserInfo)
		selfRoutes.PUT("/info", allHandlers.userHdl.EditSelfInfo)
	}

	// Swagger documentation route
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
