package controller

import (
	"net/http"
	"template/app"
	"template/dto"
	"template/model"
	"template/service"
	"template/util/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthCtn struct {
	auth   service.IAuthService
	logger *zap.SugaredLogger
}

func NewAuthCtn() app.Controller {
	var controller *AuthCtn

	// Use the mock service for testing
	app.Invoke(func(loginService service.IAuthService, logger *zap.SugaredLogger) {
		// create controller
		controller = &AuthCtn{
			auth:   loginService,
			logger: logger,
		}
	})

	return controller
}

func (c *AuthCtn) RegisterEndpoints(api *gin.RouterGroup) {
	// create a group with the name of the router
	group := api.Group("/auth")

	// register Endpoints
	group.POST("/login", c.login)
	group.POST("/refresh", c.RefreshToken)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticates a user and returns access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginDto	body		dto.LoginDto	true	"Login credentials"
//	@Success		200			{object}	dto.TokenDto
//	@Router			/auth/login [post]
func (l *AuthCtn) login(c *gin.Context) {
	var loginDto dto.LoginDto

	if err := c.BindJSON(&loginDto); err != nil {
		l.logger.Errorf("Invalid login request err = %+v", err)
		return
	}

	accessToken, refreshToken, err := l.auth.Login(loginDto.Email, loginDto.Password)
	if err != nil {
		l.logger.Errorf("Login failed err = %+v", err)
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	zap.S().Debugf("Refresh: %s", refreshToken)

	c.JSON(http.StatusOK, dto.TokenDto{
		AccessToken: accessToken,
	})
}

// Refresh godoc
//
//	@Summary		Refresh Access Token
//	@Description	Generates a new access token using a valid refresh token
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	dto.TokenDto
//	@Router			/auth/refresh [post]
func (l *AuthCtn) RefreshToken(c *gin.Context) {
	_, claims, err := auth.ParseToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		l.logger.Errorf("Error Parsing clames err = %+v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	userUuid, err := uuid.Parse(claims.ID)
	if err != nil {
		l.logger.Errorf("Error Parsing uuid err = %+v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, refreshNew, err := l.auth.RefreshTokens(&model.User{
		Uuid:  userUuid,
		Email: claims.Email,
		Role:  claims.Role,
	})
	if err != nil {
		l.logger.Error("Refresh failed err = %+v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	zap.S().Debugf("Refresh: %s", refreshNew)

	c.JSON(http.StatusOK, dto.TokenDto{
		AccessToken: token,
	})
}
