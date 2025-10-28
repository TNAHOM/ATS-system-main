package handler

import "github.com/gin-gonic/gin"

type User interface {
	SignUp(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
}

type JobPost interface {
	CreateJobPost(ctx *gin.Context)
	GetAllJobPostsByUserId(ctx *gin.Context)
	GetAllJobPosts(ctx *gin.Context)
	GetJobPostByID(ctx *gin.Context)
	UpdateJobPost(ctx *gin.Context)
	DeleteJobPost(ctx *gin.Context)
}

type Application interface {
	// GetApplicationsByApplicantID(ctx *gin.Context)
	GetApplicationsByJobPostID(ctx *gin.Context)
}
