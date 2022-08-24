package jwtpack

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	. "github.com/kadirgonen/movie-api/api/model"
	. "github.com/kadirgonen/movie-api/app/models"
	"github.com/kadirgonen/movie-api/app/pkg/config"
	. "github.com/kadirgonen/movie-api/app/pkg/errors"
	"go.uber.org/zap"
)

type TokenDetails struct {
	AccessTokenDetails  AccessTokenDetails
	RefreshTokenDetails RefreshTokenDetails
	jwt.StandardClaims
}

type AccessTokenDetails struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	Role        bool   `json:"role"`
	AccessToken string `json:"access_token"`
	AccessUuid  string `json:"access_uuid"`
	AtExpires   int64  `json:"exp"`
	jwt.StandardClaims
}
type RefreshTokenDetails struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Role         bool   `json:"role"`
	RefreshToken string `json:"refresh_token"`
	RefreshUuid  string `json:"refresh_uuid"`
	RtExpires    int64  `json:"exp"`
	jwt.StandardClaims
}

func GenerateToken(user *User, cfg *config.Config) (*Token, error) {
	td := &TokenDetails{}
	tkn := &Token{}

	// AccessToken Create
	td.AccessTokenDetails.AtExpires = time.Now().Add(time.Duration(cfg.JWTConfig.AccessSessionTime) * time.Second).Unix() // 15 minute
	td.AccessTokenDetails.AccessUuid = uuid.New().String()
	td.AccessTokenDetails.UserID = user.Id
	td.AccessTokenDetails.Email = user.Email
	td.AccessTokenDetails.Role = user.IsAdmin

	// RefreshToken Create
	td.RefreshTokenDetails.UserID = user.Id
	td.RefreshTokenDetails.Email = user.Email
	td.RefreshTokenDetails.RtExpires = time.Now().Add(time.Duration(cfg.JWTConfig.RefreshSessionTime) * time.Second).Unix() // 7 days
	td.RefreshTokenDetails.RefreshUuid = uuid.New().String()
	td.RefreshTokenDetails.Role = user.IsAdmin

	// AccessToken Claim
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessTokenDetails.AccessUuid
	atClaims["user_id"] = td.AccessTokenDetails.UserID
	atClaims["email"] = td.AccessTokenDetails.Email
	atClaims["role"] = td.AccessTokenDetails.Role
	atClaims["exp"] = td.AccessTokenDetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessTokenResult, err := at.SignedString([]byte(cfg.JWTConfig.SecretKey))
	if err != nil {
		zap.L().Debug("generatetoken.accesstoken: ", zap.Error(err))
		return nil, NewRestError(http.StatusBadRequest, "Issue on generating access token", nil)
	}

	// Token Details
	td.AccessTokenDetails.AccessToken = accessTokenResult

	// RefreshToken Claim
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshTokenDetails.RefreshUuid
	rtClaims["user_id"] = td.RefreshTokenDetails.UserID
	rtClaims["email"] = td.RefreshTokenDetails.Email
	rtClaims["role"] = td.RefreshTokenDetails.Role
	rtClaims["exp"] = td.RefreshTokenDetails.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshTokenDetails.RefreshToken, err = rt.SignedString([]byte(cfg.JWTConfig.SecretKey))
	if err != nil {
		zap.L().Debug("generatetoken.refreshtoken: ", zap.Error(err))
		return nil, NewRestError(http.StatusBadRequest, "Issue on generating refresh token", nil)
	}

	tkn.AccessToken = td.AccessTokenDetails.AccessToken
	tkn.RefreshToken = td.RefreshTokenDetails.RefreshToken
	return tkn, nil
}

// VerifyACToken returns access token if it is not expired
func VerifyACToken(token *Token, cfg *config.Config) (*AccessTokenDetails, error) {
	hmacSecretString := cfg.JWTConfig.SecretKey
	hmacSecretAc := []byte(hmacSecretString)

	accesstoken, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return hmacSecretAc, nil
	})
	if err != nil {
		zap.L().Debug("jwt.verifytokenParse: ", zap.Error(err))
		return nil, err
	}

	if !accesstoken.Valid {
		zap.L().Debug("jwt.verifytokenValid: AccessToken is not valid!")
		return nil, NewRestError(http.StatusBadRequest, "Access token is not valid!", nil)
	}

	decodedClaims := accesstoken.Claims.(jwt.MapClaims)

	var accessTokenDetails AccessTokenDetails

	jsonStringAT, _ := json.Marshal(decodedClaims)
	json.Unmarshal(jsonStringAT, &accessTokenDetails)

	accessTokenDetails.AccessToken = token.AccessToken

	return &accessTokenDetails, nil

}

// VerifyRFToken returns refresh token if it is not expired.
func VerifyRFToken(token *Token, cfg *config.Config) (*RefreshTokenDetails, error) {
	hmacSecretString := cfg.JWTConfig.SecretKey
	hmacSecretRf := []byte(hmacSecretString)
	refreshtoken, err := jwt.Parse(token.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return hmacSecretRf, nil
	})
	if err != nil {
		zap.L().Debug("jwt.verifytokenParse: ", zap.Error(err))
		return nil, err
	}

	if !refreshtoken.Valid {
		zap.L().Debug("jwt.verifytokenValid: RefreshToken is not valid!")
		return nil, NewRestError(http.StatusBadRequest, "Refresh token is not valid!", nil)
	}

	decodedClaimsRT := refreshtoken.Claims.(jwt.MapClaims)

	var refreshTokenDetails RefreshTokenDetails

	jsonStringRef, _ := json.Marshal(decodedClaimsRT)
	json.Unmarshal(jsonStringRef, &refreshTokenDetails)

	refreshTokenDetails.RefreshToken = token.RefreshToken

	return &refreshTokenDetails, nil
}

// VerifyTokenMiddleware helps the middleware for authorization
func VerifyTokenMiddleware(token string, cfg *config.Config) (*AccessTokenDetails, error) {
	hmacSecretString := cfg.JWTConfig.SecretKey
	hmacSecretRf := []byte(hmacSecretString)
	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return hmacSecretRf, nil
	})

	if err != nil {
		zap.L().Debug("jwt.verifytokenMiddleware: ", zap.Error(err))
		return nil, err
	}

	if !parsedtoken.Valid {
		zap.L().Debug("jwt.verifytoken.middleware.Valid: Token is not valid!")
		return nil, NewRestError(http.StatusBadRequest, "Token is not valid!", nil)
	}

	decodedClaims := parsedtoken.Claims.(jwt.MapClaims)

	var tokenDetails AccessTokenDetails

	jsonStringRef, _ := json.Marshal(decodedClaims)
	json.Unmarshal(jsonStringRef, &tokenDetails)
	return &tokenDetails, nil
}
