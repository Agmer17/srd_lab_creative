package auth

import (
	"net/http"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

const oneWeek = 7 * 24 * 60 * 60

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(svc *AuthService) *AuthHandler {
	return &AuthHandler{
		service: svc,
	}
}

func (ah *AuthHandler) HandleGoogleLogin(c *gin.Context) {
	redUrl, err := ah.service.GetLoginGoogleUrl()
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, redUrl)
}

func (ah *AuthHandler) HandleGoogleCallback(c *gin.Context) {

	code := c.Query("code")

	accessToken, refreshToken, err := ah.service.AuthenticateGoogleUser(c.Request.Context(), code)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		oneWeek,
		"/",
		"localhost",
		false,
		true,
	)

	// GANTI INI SAMA REDIRECT FE NANTINYA!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	c.JSON(200, shared.SuccessResponse{Code: 200, Message: "successfully login with google", Data: accessToken})
}

func (ah *AuthHandler) LogoutHandler(c *gin.Context) {

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"message": "logged out",
	})
}

func (ah *AuthHandler) HandleRefreshSession(c *gin.Context) {

	userId, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(401, "no cookies found! can't refresh the current session!")
	}

	accessToken, err := ah.service.GetRefreshSession(c.Request.Context(), userId)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.SuccessResponse{Code: 200, Message: "refreshing the session successfully", Data: accessToken})
}

func (ah *AuthHandler) RegisterRoutes(r gin.IRouter) {
	auth := r.Group("/auth")
	{
		auth.GET("/login/google", ah.HandleGoogleLogin)
		auth.GET("/google-callback", ah.HandleGoogleCallback)
	}

	privateAuth := r.Group("/auth")
	{
		privateAuth.GET("/logout", ah.LogoutHandler)
		privateAuth.GET("/refresh-session", ah.HandleRefreshSession)
	}
}
