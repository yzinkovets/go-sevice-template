package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-service-template/config"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const TokenPayloadKey = "tokenPayload"

type JwtAuth struct {
	jwks *keyfunc.JWKS
	cfg  *config.JwtAuthConfig
}

func NewJwtAuth(cfg *config.JwtAuthConfig) (*JwtAuth, error) {
	if cfg.JwksFilePath == "" {
		return nil, errors.New("jwks file settings empty")
	}

	dat, err := os.ReadFile(cfg.JwksFilePath)
	if err != nil {
		return nil, fmt.Errorf("JWKS file read error: %s", err.Error())
	}

	var jwksJSON = json.RawMessage(dat)
	jwks, err := keyfunc.NewJSON(jwksJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS from JSON error: %s", err.Error())
	}

	return &JwtAuth{jwks: jwks, cfg: cfg}, nil
}

func (j *JwtAuth) Handle(ctx fiber.Ctx) error {

	authHeader := ctx.Get("Authorization", "")

	if authHeader == "" {
		logrus.Errorln("empty Authorization")
		return fiber.ErrUnauthorized
	}

	reqToken := strings.ReplaceAll(authHeader, "Bearer ", "")

	// do cheap offline token validation
	tokenPayload, ok := jwtTokenValidate(reqToken, j.jwks)
	if !ok || time.Now().Unix() > tokenPayload.Exp {
		logrus.Error("Wrong token or outdated")
		return fiber.ErrUnauthorized
	}

	logrus.Trace("Check token in external service: ", j.cfg.JwtCheckUrl, " timeout: ", j.cfg.JwtCheckUrlTimeoutSec, " sec..")

	// go to external service to check token
	if err := j.checkTokenWithKeycloakService(authHeader); err != nil {
		return err
	}

	// everything is ok, set token payload to context
	ctx.Context().SetUserValue(TokenPayloadKey, tokenPayload)

	return ctx.Next()
}

func (j *JwtAuth) checkTokenWithKeycloakService(authHeader string) *fiber.Error {
	client := http.Client{
		Timeout: time.Duration(j.cfg.JwtCheckUrlTimeoutSec) * time.Second,
	}

	request, err := http.NewRequest(http.MethodHead, j.cfg.JwtCheckUrl, nil)
	if err != nil {
		logrus.Error("Error while checking token in external service: ", err)
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	request.Header.Add("Authorization", authHeader)
	resp, err := client.Do(request)
	if err != nil {
		logrus.Error("Error while checking token in external service: ", err)
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		errMessage := fmt.Sprintf("Error while checking token in external service. StatusCode: %d", resp.StatusCode)
		logrus.Error(errMessage)
		return fiber.NewError(resp.StatusCode, errMessage)
	}

	return nil
}

func jwtTokenValidate(tokenString string, jwks *keyfunc.JWKS) (*TokenPayloadFields, bool) {
	payload := &TokenPayloadFields{}
	token, err := jwt.ParseWithClaims(tokenString, payload, jwks.Keyfunc)
	if err != nil {
		// log.Fatal(err)
		return nil, false
	}

	return payload, token.Valid
}

func GetTokenPayload(c fiber.Ctx) (*TokenPayloadFields, error) {
	tokenPayload, converted := c.Context().UserValue(TokenPayloadKey).(*TokenPayloadFields)

	if !converted {
		return nil, errors.New("cannot find token payload")
	}

	return tokenPayload, nil
}

func GetTokenEmail(c fiber.Ctx) (string, error) {

	tokenPayload, err := GetTokenPayload(c)

	if err != nil {
		return "", err
	}

	return tokenPayload.Email, nil
}

func GetReqEmail(paramEmail string, c fiber.Ctx) (string, error) {

	tokenPayload, err := GetTokenPayload(c)

	if err != nil {
		// if we don't have token (from router)
		if paramEmail != "" {
			return paramEmail, nil
		}

		return "", err
	}

	email := paramEmail

	if paramEmail != "" {
		if tokenPayload.Email != paramEmail && !tokenPayload.IsAdmin() && !tokenPayload.IsSupport() {
			return "", errors.New("wrong email param")
		}
	} else {
		email = tokenPayload.Email
	}

	return email, nil
}

type TokenPayloadFields struct {
	Exp          int64  `json:"exp"`
	Iss          string `json:"iss"`
	Azp          string `json:"azp"`
	Sub          string `json:"sub"` // user_id
	Email        string `json:"email"`
	Name         string `json:"name"`
	SessionState string `json:"session_state"`
	IsAdminInt   int    `json:"custom:isadmin"`
	IsSupportInt int    `json:"custom:issupport"`
	jwt.RegisteredClaims
}

func (there *TokenPayloadFields) IsAdmin() bool {
	return there.IsAdminInt == 1
}

func (there *TokenPayloadFields) IsSupport() bool {
	return there.IsSupportInt == 1
}
