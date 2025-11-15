package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/salex06/pr-service/internal/config"
	"github.com/salex06/pr-service/internal/database"
	prRepo "github.com/salex06/pr-service/internal/repos/pr"
	revsRepo "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepo "github.com/salex06/pr-service/internal/repos/team"
	userRepo "github.com/salex06/pr-service/internal/repos/user"
	"github.com/salex06/pr-service/internal/rest"
	"github.com/salex06/pr-service/internal/service"
)

func main() {
	dbConfig := config.Load()
	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Println(fmt.Errorf("connecting to database failed: %w", err))
		return
	}

	var tr teamRepo.TeamRepository = teamRepo.NewPostgresTeamRepository(db)
	var ur userRepo.UserRepository = userRepo.NewPostgresUserRepository(db)
	var ar revsRepo.AssignedRevsRepository = revsRepo.NewPostgresAssignedRevsRepository(db)
	var pr prRepo.PullRequestRepository = prRepo.NewPostgresPullRequestRepository(db)

	var ts *service.TeamService = service.NewTeamService(&tr, &ur)
	var us *service.UserService = service.NewUserService(&ur, &ar, &pr)
	var prs *service.PullRequestService = service.NewPullRequestService(&pr, &ar, &ur, &tr)

	var th *rest.TeamHandler = rest.NewTeamHandler(ts)
	var uh *rest.UserHandler = rest.NewUserHandler(us)
	var prh *rest.PullRequestHandler = rest.NewPullRequestHandler(prs)

	r := gin.Default()

	setupTeamHandlers(th, r)
	setupUserHandlers(uh, r)
	setupPullRequestHandlers(prh, r)

	r.Run(":8080")
}

func setupTeamHandlers(th *rest.TeamHandler, r *gin.Engine) {
	r.POST("/team/add", th.HandleAddTeamRequest)
	r.GET("/team/get", th.HandleGetTeamRequest)
}

func setupUserHandlers(uh *rest.UserHandler, r *gin.Engine) {
	r.POST("/users/setIsActive", uh.HandleSetIsActiveRequest)
	r.GET("/users/getReview", uh.HandleGetReviewRequest)
}

func setupPullRequestHandlers(prh *rest.PullRequestHandler, r *gin.Engine) {
	r.POST("/pullRequest/create", prh.HandleCreateRequest)
	r.POST("/pullRequest/merge", prh.HandleMergeRequest)
	r.POST("/pullRequest/reassign", prh.HandleReassignRequest)
}
