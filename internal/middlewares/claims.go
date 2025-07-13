package middlewares

import (
	"errors"

	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/labstack/echo/v4"
)

const (
	userClaimsKey = "user"
	otpClaimsKey  = "otp"
)

var (
	ErrUserClaimsNotFound = errors.New("user claims not found")
	ErrOTPClaimsNotFound  = errors.New("otp claims not found")
)

func SetUserClaims(ectx echo.Context, claims *models.AccessTokenClaims) {
	ectx.Set(userClaimsKey, claims)
}

func SetOTPClaims(ectx echo.Context, claims *models.OTPTokenClaims) {
	ectx.Set(otpClaimsKey, claims)
}

func GetUserClaims(ectx echo.Context) (*models.AccessTokenClaims, error) {
	user, ok := ectx.Get(userClaimsKey).(*models.AccessTokenClaims)
	if !ok {
		return nil, ErrUserClaimsNotFound
	}

	return user, nil
}

func GetOTPClaims(ectx echo.Context) (*models.OTPTokenClaims, error) {
	otp, ok := ectx.Get(otpClaimsKey).(*models.OTPTokenClaims)
	if !ok {
		return nil, ErrOTPClaimsNotFound
	}

	return otp, nil
}

func GetOTPJTI(ectx echo.Context) string {
	otp, err := GetOTPClaims(ectx)
	if err != nil {
		return ""
	}

	return otp.ID
}

func GetUserID(ectx echo.Context) string {
	user, err := GetUserClaims(ectx)
	if err != nil {
		return ""
	}

	return user.Subject
}
