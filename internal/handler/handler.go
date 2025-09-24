package handler

import "github.com/gin-gonic/gin"

type User interface {
	SignUp(ctx *gin.Context)
	LoginUser(ctx *gin.Context)

	GetAllUsers(ctx *gin.Context)
}

type JobPost interface {
	CreateJobPost(ctx *gin.Context)
	GetAllJobPosts(ctx *gin.Context)
}
