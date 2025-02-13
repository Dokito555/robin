package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewAuth(userService *services.UserService, tokenService *services.TokenService) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        tokenStr := ctx.GetHeader("Authorization")
        
        request := &model.VerifyUserRequest{Token: tokenStr}
        
        _, err := userService.Verify(ctx.Request.Context(), request)
        if err != nil {
            userService.Log.Warnf("[Auth Middleware] Failed to find user by token: %+v", err)
            ctx.JSON(http.StatusUnauthorized, nil)
            ctx.Abort()
            return 
        }

        claim, err := tokenService.ValidateToken(ctx.Request.Context(), tokenStr)
        if err != nil {
            userService.Log.Warnf("failed to validate token: %+v", err)
            ctx.JSON(http.StatusUnauthorized, nil)
            ctx.Abort()
            return
        }

        if time.Now().Unix() > claim.ExpiresAt.Unix() {
            userService.Log.Warnf("JWT token is expired. Expiry: %v, Current: %v", claim.ExpiresAt, time.Now())
            ctx.JSON(http.StatusUnauthorized, nil)
            ctx.Abort()
            return 
        }

        ctx.Set("auth", claim)
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