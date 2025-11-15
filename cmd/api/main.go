package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/salex06/pr-service/internal/config"
	"github.com/salex06/pr-service/internal/database"
	prRepository "github.com/salex06/pr-service/internal/repos/pr"
	revsRepository "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepository "github.com/salex06/pr-service/internal/repos/team"
	userRepository "github.com/salex06/pr-service/internal/repos/user"
	"github.com/salex06/pr-service/internal/rest"
	"github.com/salex06/pr-service/internal/service"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	dbConfig := config.LoadDBConfig()
	appConfig := config.LoadAppConfig()

	// Подключение к БД
	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Println(fmt.Errorf("connecting to database failed: %w", err))
		return
	}

	// Инициализация и внедрение компонентов приложения
	teamRepo := teamRepository.NewPostgresTeamRepository(db)
	userRepo := userRepository.NewPostgresUserRepository(db)
	revsRepo := revsRepository.NewPostgresAssignedRevsRepository(db)
	pullRequestRepo := prRepository.NewPostgresPullRequestRepository(db)

	teamService := service.NewTeamService(&teamRepo, &userRepo)
	userService := service.NewUserService(&userRepo, &revsRepo, &pullRequestRepo)
	pullRequestService := service.NewPullRequestService(&pullRequestRepo, &revsRepo, &userRepo, &teamRepo)
	statService := service.NewStatsService(&pullRequestRepo, &revsRepo, &userRepo, &teamRepo)

	teamHandler := rest.NewTeamHandler(teamService)
	userHandler := rest.NewUserHandler(userService)
	pullRequestHandler := rest.NewPullRequestHandler(pullRequestService)
	statsHandler := rest.NewStatHandler(statService)

	r := gin.Default()

	// Настройка эндпоинтов
	setupTeamHandlers(teamHandler, r)
	setupUserHandlers(userHandler, r)
	setupPullRequestHandlers(pullRequestHandler, r)
	setupStatRequestHandlers(statsHandler, r)

	// Запуск сервера
	err = r.Run(fmt.Sprintf(":%s", appConfig.ServerPort))
	if err != nil {
		log.Printf("unable to start server: %s\n", err)
	}
}

func setupTeamHandlers(handler *rest.TeamHandler, r *gin.Engine) {
	r.POST("/team/add", handler.HandleAddTeamRequest)
	r.GET("/team/get", handler.HandleGetTeamRequest)
}

func setupUserHandlers(handler *rest.UserHandler, r *gin.Engine) {
	r.POST("/users/setIsActive", handler.HandleSetIsActiveRequest)
	r.GET("/users/getReview", handler.HandleGetReviewRequest)
}

func setupPullRequestHandlers(handler *rest.PullRequestHandler, r *gin.Engine) {
	r.POST("/pullRequest/create", handler.HandleCreateRequest)
	r.POST("/pullRequest/merge", handler.HandleMergeRequest)
	r.POST("/pullRequest/reassign", handler.HandleReassignRequest)
}

func setupStatRequestHandlers(handler *rest.StatsHandler, r *gin.Engine) {
	r.GET("/stats", handler.HandleGetStatsRequest)
}
