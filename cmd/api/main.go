package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/salex06/pr-service/internal/config"
	"github.com/salex06/pr-service/internal/database"
	prRepo "github.com/salex06/pr-service/internal/repos/pr"
	revsRepo "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepo "github.com/salex06/pr-service/internal/repos/team"
	userRepo "github.com/salex06/pr-service/internal/repos/user"
	"github.com/salex06/pr-service/internal/rest"
	"github.com/salex06/pr-service/internal/service"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	dbConfig := config.LoadDbConfig()
	appConfig := config.LoadAppConfig()

	//Подключение к БД
	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Println(fmt.Errorf("connecting to database failed: %w", err))
		return
	}

	//Инициализация и внедрение компонентов приложения
	teamRepo := teamRepo.NewPostgresTeamRepository(db)
	userRepo := userRepo.NewPostgresUserRepository(db)
	revsRepo := revsRepo.NewPostgresAssignedRevsRepository(db)
	pullRequestRepo := prRepo.NewPostgresPullRequestRepository(db)

	teamService := service.NewTeamService(&teamRepo, &userRepo)
	userService := service.NewUserService(&userRepo, &revsRepo, &pullRequestRepo)
	pullRequestService := service.NewPullRequestService(&pullRequestRepo, &revsRepo, &userRepo, &teamRepo)

	teamHandler := rest.NewTeamHandler(teamService)
	userHandler := rest.NewUserHandler(userService)
	pullRequestHandler := rest.NewPullRequestHandler(pullRequestService)

	r := gin.Default()

	//Настройка эндпоинтов
	setupTeamHandlers(teamHandler, r)
	setupUserHandlers(userHandler, r)
	setupPullRequestHandlers(pullRequestHandler, r)

	//Запуск сервера
	r.Run(fmt.Sprintf(":%s", appConfig.ServerPort))
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
