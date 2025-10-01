package main

import (
	"log"

	"github.com/gin-gonic/gin"
//	"github.com/gin-contrib/cors" 
	"github.com/saku-730/specimen-web/backend/config"
	"github.com/saku-730/specimen-web/backend/internal/handler"
	"github.com/saku-730/specimen-web/backend/internal/infrastructure"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"github.com/saku-730/specimen-web/backend/internal/service"
)

func main() {
	// 1. 設定を読み込む (ステップ1)
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗: %v", err)
	}

	// 2. データベースに接続する (ステップ2)
	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("データベース接続に失敗: %v", err)
	}

	// --- ここからが「依存性の注入」による組み立て作業！ ---

	// 3. repository層を初期化 (ステップ4)
	userRepo := repository.NewUserRepository(db)
	// (他のprojectRepoなどもここで初期化する)

	// 4. service層を初期化 (ステップ5)
	// userServiceは、userRepoとdb接続に依存している
	userService := service.NewUserService(db, userRepo)
	// (他のprojectServiceなどもここで初期化する)

	// userHandlerは、userServiceに依存している
	userHandler := handler.NewUserHandler(userService)
	// (他のprojectHandlerなどもここで初期化する)

	// --- 組み立て完了！ ---

	// 6. ルーターをセットアップする
	router := gin.Default()
	apiV1 := router.Group("/api/v0-0-1") // APIのバージョニング
	{
		// userHandlerに、ルーターへのルート登録を依頼する
		userHandler.RegisterUserRoutes(apiV1)
		// (他のhandlerのルートもここで登録する)
	}

	// 7. サーバーを起動する！
	log.Printf("サーバーをポート %s で起動します...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("サーバーの起動に失敗: %v", err)
	}
}
