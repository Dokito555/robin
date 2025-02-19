package middleware

import (
	"log"
	"net/http"

	"github.com/Dokito555/robin-songs/internal/model"
	"github.com/Dokito555/robin-songs/internal/services"
	"github.com/Dokito555/robin-songs/pkg/artist"
	"github.com/Dokito555/robin-songs/utils/constants"
	"github.com/gin-gonic/gin"
)

func NewAuth(artist *artist.ArtistPkg, tokenService *services.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")

		// profile, err := artist.GetProfile(ctx.Request.Context(), tokenStr)
		profile, err := artist.ValidateToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			tokenService.Log.Warnf("failed find user by token : %+v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			ctx.Abort()
			return
		}

		claimToken := &model.ClaimToken{
			UserID:   profile.ID,
			Email:    profile.Email,
			Role:     profile.Role,
			UserName: profile.UserName,
		}

		if claimToken.UserID == 0 || claimToken.Email == "" || claimToken.Role == "" {
			tokenService.Log.WithField("profile", claimToken).Warn("Invalid user profile received from artist service")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user profile"})
			ctx.Abort()
			return
		}

		if claimToken.Role != constants.ROLE_ADMIN && claimToken.Role != constants.ROLE_ARTIST {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("auth", claimToken)
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