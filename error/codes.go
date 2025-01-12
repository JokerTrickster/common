package error

/*
	에러 코드 정의 및 매핑
*/

import "net/http"

// IErrFrom defines where the error originated
type IErrFrom string

const (
	ErrFromClient   = IErrFrom("client")
	ErrFromInternal = IErrFrom("internal")
	ErrFromMongoDB  = IErrFrom("mongoDB")
	ErrFromMysqlDB  = IErrFrom("mysqlDB")
	ErrFromRedis    = IErrFrom("redis")
	ErrFromAws      = IErrFrom("aws")
	ErrFromAwsS3    = IErrFrom("aws_s3")
	ErrFromAwsSsm   = IErrFrom("aws_ssm")
	ErrFromNaver    = IErrFrom("naver")
	ErrFromGemini   = IErrFrom("gemini")
	ErrFromKakao    = IErrFrom("kakao")
	ErrFromFirebase = IErrFrom("firebase")
	ErrFromChatGPT  = IErrFrom("chatGPT")
)

// ErrType defines the type of errors
type ErrType string

const (
	ErrBadParameter   = ErrType("PARAM_BAD")
	ErrNotFound       = ErrType("NOT_FOUND")
	ErrBadToken       = ErrType("TOKEN_BAD")
	ErrInternalServer = ErrType("INTERNAL_SERVER")
	ErrInternalDB     = ErrType("INTERNAL_DB")
	ErrPartner        = ErrType("PARTNER")

	// Auth errors
	ErrCodeNotFound           = ErrType("CODE_NOT_FOUND")
	ErrUserNotFound           = ErrType("USER_NOT_FOUND")
	ErrProfileNotFount        = ErrType("PROFILE_NOT_FOUND")
	ErrUserAlreadyExisted     = ErrType("USER_ALREADY_EXISTED")
	ErrInvalidAccessToken     = ErrType("INVALID_ACCESS_TOKEN")
	ErrPasswordNotMatch       = ErrType("PASSWORD_NOT_MATCH")
	ErrInvalidAuthCode        = ErrType("INVALID_AUTH_CODE")
	ErrInvalidEmailOrPassword = ErrType("INVALID_EMAIL_OR_PASSWORD")

	// Food errors
	ErrGeminiError  = ErrType("GEMINI_INTERNAL_SERVER")
	ErrFoodNotFound = ErrType("FOOD_NOT_FOUND")
)

// ErrHttpCode maps error types to HTTP status codes
var ErrHttpCode = map[string]int{
	"PARAM_BAD":                 http.StatusBadRequest,
	"USER_ALREADY_EXISTED":      http.StatusBadRequest,
	"BAD_REQUEST":               http.StatusBadRequest,
	"USER_NOT_FOUND":            http.StatusBadRequest,
	"NOT_ENOUGH_CARD":           http.StatusBadRequest,
	"NOT_ENOUGH_CONDITION":      http.StatusBadRequest,
	"PASSWORD_NOT_MATCH":        http.StatusBadRequest,
	"INVALID_EMAIL_OR_PASSWORD": http.StatusBadRequest,
	"FOOD_NOT_FOUND":            http.StatusBadRequest,
	"TOKEN_BAD":                 http.StatusUnauthorized,
	"INVALID_ACCESS_TOKEN":      http.StatusUnauthorized,
	"INVALID_AUTH_CODE":         http.StatusUnauthorized,
	"PARTNER":                   http.StatusForbidden,
	"NOT_FOUND":                 http.StatusNotFound,
	"INTERNAL_SERVER":           http.StatusInternalServerError,
	"INTERNAL_DB":               http.StatusInternalServerError,
	"GEMINI_INTERNAL_SERVER":    http.StatusInternalServerError,
}
