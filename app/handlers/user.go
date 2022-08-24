package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	. "github.com/kadirgonen/movie-api/api/model"
	. "github.com/kadirgonen/movie-api/app/models"
	. "github.com/kadirgonen/movie-api/app/pkg/errors"
	. "github.com/kadirgonen/movie-api/app/service"

	"github.com/kadirgonen/movie-api/app/pkg/config"
	. "github.com/kadirgonen/movie-api/app/pkg/helper"
	. "github.com/kadirgonen/movie-api/app/pkg/jwt"

	"go.uber.org/zap"
)

type UserHandler struct {
	userService UserServiceInterface
	cfg         *config.Config
}

func NewUserHandler(r *gin.RouterGroup, u *UserService, cfg *config.Config) {
	h := &UserHandler{userService: u, cfg: cfg}
	u.Migrate()
	r.POST("/signup", h.signup)
	r.POST("/login", h.login)
}

// signup helps user to signup to system. It check user's e-mail & password validation regarding of rules. If it is completed with success, it creates Access & Refresh token and store it as a cookie.
func (u *UserHandler) signup(c *gin.Context) {
	var req SignUp
	if err := c.Bind(&req); err != nil {
		zap.L().Error("user.handler.signup", zap.Error(err))
		c.JSON(ErrorResponse(NewRestError(http.StatusBadRequest, os.Getenv("CHECK_YOUR_REQUEST"), nil)))
		return
	}

	if err := req.Validate(strfmt.NewFormats()); err != nil {
		zap.L().Error("user.handler.signup", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	res, err := u.userService.CheckUser(ResponseToUser(&req))
	if err != nil {
		zap.L().Error("user.handler.signup", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	if res {
		zap.L().Error("user.handler.signup: User Already exist")
		c.JSON(ErrorResponse(NewRestError(http.StatusBadRequest, os.Getenv("USER_ALREADY_EXIST"), nil)))
		return
	}

	user, err := u.userService.Save(ResponseToUser(&req))
	if err != nil {
		zap.L().Error("user.handler.signup", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}

	tkn, err := GenerateToken(user, u.cfg)
	if err != nil {
		zap.L().Error("user.handler.signup: generatetoken", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}

	cookie := SetCookie(tkn, user)
	if cookie != nil {
		http.SetCookie(c.Writer, cookie)
		c.JSON(http.StatusCreated, APIResponseSignUp{Code: http.StatusCreated, Token: tkn})
		return
	}

}

// login helps user to enter system. It checks user info and cookies. If it is complete in a success, return a valid token
func (u *UserHandler) login(c *gin.Context) {
	var req Login
	if err := c.Bind(&req); err != nil {
		c.JSON(ErrorResponse(NewRestError(http.StatusBadRequest, "Check your request body", nil)))
	}
	if err := req.Validate(strfmt.NewFormats()); err != nil {
		zap.L().Error("user.handler.login", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}

	user, err := u.userService.Login(*req.Email, *req.Password)
	if err != nil {
		c.JSON(ErrorResponse(err))
		return
	}

	value, err := DecodeCookie(c.Request, user)
	if err != nil {
		zap.L().Error("user.handler.login: decodetoken", zap.Error(err))
		c.JSON(ErrorResponse(err))
		return
	}
	tokendetails, err := VerifyACToken(value, u.cfg)
	if err != nil {
		zap.L().Error("user.handler.login: verifyactoken", zap.Error(err))
		if err.Error() == "Token is expired" {
			rftokendetails, err := VerifyRFToken(value, u.cfg)
			if err != nil {
				zap.L().Error("user.handler.login: verifyrftoken", zap.Error(err))
				if err.Error() == "Token is expired" {
					tkn, err := GenerateToken(user, u.cfg)
					if err != nil {
						zap.L().Error("user.handler.signup: generatetoken", zap.Error(err))
						c.JSON(ErrorResponse(err))
						return
					}
					cookie := SetCookie(tkn, user)
					if cookie != nil {
						http.SetCookie(c.Writer, cookie)
						c.JSON(http.StatusCreated, APIResponseSignUp{Code: http.StatusCreated, Token: tkn})
						return
					}
				}
				c.JSON(ErrorResponse(err))
				return
			}
			if rftokendetails.UserID == user.Id {
				c.JSON(http.StatusOK, SoleToken{Code: http.StatusOK, Token: rftokendetails.RefreshToken})
				return
			}

		}
		c.JSON(ErrorResponse(err))
		return
	}

	if tokendetails.UserID == user.Id {
		c.JSON(http.StatusOK, SoleToken{Code: http.StatusOK, Token: tokendetails.AccessToken})
		return
	}

}

func (u *UserHandler) Migrate() {
	u.userService.Migrate()
}
