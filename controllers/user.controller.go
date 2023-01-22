package controllers

import (
	"github.com/SarathLUN/auth-service-grpc-golang/models"
	"github.com/SarathLUN/auth-service-grpc-golang/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
