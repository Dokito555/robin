package middleware

import (
	"log"
	"net/http"

	model "github.com/Dokito555/robin-albums/internal/models"
	"github.com/Dokito555/robin-albums/internal/services"
	ums "github.com/Dokito555/robin-albums/pkg/ums"
	"github.com/Dokito555/robin-albums/utils/constants"
	"github.com/gin-gonic/gin"
)

func NewAuth(ums *ums.UMS, tokenService *services.TokenService) gin.HandlerFunc {
    return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")

		// profile, err := ums.GetProfile(ctx.Request.Context(), tokenStr)
		profile, err := ums.ValidateToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			tokenService.Log.Warnf("failed find user by token : %+v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			ctx.Abort()
			return
		}

		if profile.ID == 0 || profile.Email == "" || profile.Role == "" {
			tokenService.Log.Warnf("Invalid user profile received from UMS: %+v", profile)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user profile"})
			ctx.Abort()
			return
		}

		if profile.Role != constants.ROLE_ADMIN && profile.Role != constants.ROLE_ARTIST {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("auth", profile)
		ctx.Next()
	}
}

func GetProfile(ctx *gin.Context) *model.ClaimToken {
    auth, exist := ctx.Get("auth")
    if !exist {
        log.Printf("auth not found in context")
        return nil
    }
    
    claim, ok := auth.(*model.ClaimToken)
    if !ok {
        log.Printf("type assertion failed. Expected *model.ClaimToken, got type: %T, value: %+v", 
            auth, auth)
        return nil
    }
    
    return claim
}