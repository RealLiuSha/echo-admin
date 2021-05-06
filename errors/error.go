package errors

import (
	"github.com/pkg/errors"
)

// Alias
// noinspection GoUnusedGlobalVariable
var (
	Is           = errors.Is
	As           = errors.As
	New          = errors.New
	Unwrap       = errors.Unwrap
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

// Database
var (
	DatabaseInternalError  = errors.New("database internal error")
	DatabaseRecordNotFound = errors.New("database record not found")
)

// Redis
var (
	RedisKeyNoExist = errors.New("redis key does not exist")
)

// Captcha
var (
	CaptchaAnswerCodeNoMatch = errors.New("captcha answer code no match")
)

// Auth
var (
	AuthTokenInvalid      = errors.New("auth token is invalid")
	AuthTokenExpired      = errors.New("auth token is expired")
	AuthTokenNotValidYet  = errors.New("auth token not active yet")
	AuthTokenMalformed    = errors.New("auth token is malformed")
	AuthTokenGenerateFail = errors.New("failed to generate auth token")
)
