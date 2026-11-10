package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/api/middleware"
	bookmarkHdl "github.com/ThienKim52/golang-dev/internal/app/handler/bookmark"
	genpassHdl "github.com/ThienKim52/golang-dev/internal/app/handler/genpass"
	healthcheckHdl "github.com/ThienKim52/golang-dev/internal/app/handler/healthcheck"
	linkHdl "github.com/ThienKim52/golang-dev/internal/app/handler/link"
	userHdl "github.com/ThienKim52/golang-dev/internal/app/handler/user"
	bookmarkRepo "github.com/ThienKim52/golang-dev/internal/app/repository/bookmark"
	repoLink "github.com/ThienKim52/golang-dev/internal/app/repository/link"
	repo "github.com/ThienKim52/golang-dev/internal/app/repository/user"
	bookmarkSvc "github.com/ThienKim52/golang-dev/internal/app/service/bookmark"
	genpassSvc "github.com/ThienKim52/golang-dev/internal/app/service/genpass"
	healthcheckSvc "github.com/ThienKim52/golang-dev/internal/app/service/healthcheck"
	linkSvc "github.com/ThienKim52/golang-dev/internal/app/service/link"
	userSvc "github.com/ThienKim52/golang-dev/internal/app/service/user"
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
	if opts == nil {
		opts = &EngineOpts{}
	}
	if opts.App == nil {
		opts.App = gin.New()
	}
	if opts.Config == nil {
		opts.Config = &Config{}
	}
	app := &engine{
		app:         opts.App,
		config:      opts.Config,
		redisClient: opts.RedisClient,
		db:          opts.DB,
		jwtGen:      opts.JWTGen,
		jwtVal:      opts.JWTVal,
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

func (e *engine) initHandlers() *allHandlers {
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

	bookmarkRepo := bookmarkRepo.NewBookmarkRepository(e.db)
	bookmarkSvc := bookmarkSvc.NewBookmarkService(bookmarkRepo, utils.NewKeyGenerator())
	bookmarkHdl := bookmarkHdl.NewBookmarkHandler(bookmarkSvc)

	return &allHandlers{
		healthCheckSvcHdl: healthCheckSvcHdl,
		genPassHdl:        genPassHdl,
		linkHdl:           linkHdl,
		userHdl:           userHdl,
		bookmarkHdl:       bookmarkHdl,
	}
}

type allHandlers struct {
	healthCheckSvcHdl healthcheckHdl.HealthCheckHandler
	genPassHdl        genpassHdl.GenPass
	linkHdl           linkHdl.LinkHandler
	userHdl           userHdl.Handler
	bookmarkHdl       bookmarkHdl.BookmarkHandler
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
	bookmarkRoutes := v1Routes.Group("/bookmarks")
	bookmarkRoutes.Use(jwtAuth.JWTAuth())
	{
		bookmarkRoutes.POST("", allHandlers.bookmarkHdl.AddBookmark)
		bookmarkRoutes.GET("", allHandlers.bookmarkHdl.GetBookmarks)
		bookmarkRoutes.PUT("/:id", allHandlers.bookmarkHdl.UpdateBookmark)
		bookmarkRoutes.DELETE("/:id", allHandlers.bookmarkHdl.DeleteBookmark)
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
