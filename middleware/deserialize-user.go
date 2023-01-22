package middleware

import (
	"fmt"
	"github.com/SarathLUN/auth-service-grpc-golang/config"
	"github.com/SarathLUN/auth-service-grpc-golang/services"
	"github.com/SarathLUN/auth-service-grpc-golang/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func DeserializeUser(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHandler := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHandler)
		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in!"})
			return
		}

		conf, _ := config.LoadConfig(".")
		sub, err := utils.ValidateToken(accessToken, conf.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		user, err := userService.FindUserById(fmt.Sprint(sub))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "The user belonging to this token no longer exists!"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}

}
