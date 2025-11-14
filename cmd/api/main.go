package main

import (
	"github.com/gin-gonic/gin"
	prRepo "github.com/salex06/pr-service/internal/repos/pr"
	revsRepo "github.com/salex06/pr-service/internal/repos/reviewers"
	teamRepo "github.com/salex06/pr-service/internal/repos/team"
	userRepo "github.com/salex06/pr-service/internal/repos/user"
	"github.com/salex06/pr-service/internal/rest"
	"github.com/salex06/pr-service/internal/service"
)

func main() {
	var tr teamRepo.TeamRepository = teamRepo.NewInMemoryTeamRepository()
	var ur userRepo.UserRepository = userRepo.NewInMemoryUserRepository()
	var ar revsRepo.AssignedRevsRepository = revsRepo.NewInMemoryAssignedRevsRepository()
	var pr prRepo.PullRequestRepository = prRepo.NewInMemoryPullRequestRepository()

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
