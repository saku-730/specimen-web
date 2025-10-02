package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/saku-730/specimen-web/backend/config"
	"github.com/saku-730/specimen-web/backend/internal/handler"
	"github.com/saku-730/specimen-web/backend/internal/infrastructure"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

func main() {
	// load config
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗: %v", err)
	}

	// connect database
	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("データベース接続に失敗: %v", err)
	}

	// --- ここからが「依存性の注入」による組み立て作業！ ---

	// Repository層を初期化
	userRepo := repository.NewUserRepository(db)
	occurrenceRepo := repository.NewOccurrenceRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	specimenRepo := repository.NewSpecimenRepository(db)
	identificationRepo := repository.NewIdentificationRepository(db)
	observationRepo := repository.NewObservationRepository(db)
	wikiRepo := repository.NewWikiRepository(db)
	_ = repository.NewAttachmentRepository(db) // service/handlerがないので一旦変数に入れない
	_ = repository.NewLogRepository(db)
	_ = repository.NewPlaceRepository(db)

	// Service層を初期化
	userService := service.NewUserService(db, userRepo)
	occurrenceService := service.NewOccurrenceService(db, occurrenceRepo)
	projectService := service.NewProjectService(db, projectRepo)
	specimenService := service.NewSpecimenService(db, specimenRepo)
	identificationService := service.NewIdentificationService(db, identificationRepo)
	observationService := service.NewObservationService(db, observationRepo)
	wikiService := service.NewWikiService(db, wikiRepo)

	// Handler層を初期化
	userHandler := handler.NewUserHandler(userService)
	occurrenceHandler := handler.NewOccurrenceHandler(occurrenceService)
	projectHandler := handler.NewProjectHandler(projectService)
	specimenHandler := handler.NewSpecimenHandler(specimenService)
	identificationHandler := handler.NewIdentificationHandler(identificationService)
	observationHandler := handler.NewObservationHandler(observationService)
	wikiHandler := handler.NewWikiHandler(wikiService)

	//setup router
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	apiV0_0_1 := router.Group("/api/v0_0_1") // APIのバージョニング
	{
		userHandler.RegisterUserRoutes(apiV0_0_1)
		occurrenceHandler.RegisterOccurrenceRoutes(apiV0_0_1)
		projectHandler.RegisterProjectRoutes(apiV0_0_1)
		specimenHandler.RegisterSpecimenRoutes(apiV0_0_1)
		identificationHandler.RegisterIdentificationRoutes(apiV0_0_1)
		observationHandler.RegisterObservationRoutes(apiV0_0_1)
		wikiHandler.RegisterWikiRoutes(apiV0_0_1)
	}

	// start server
	log.Printf("サーバーをポート %s で起動します...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("サーバーの起動に失敗: %v", err)
	}
}
