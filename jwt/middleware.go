package jwt

import (
	_error "github.com/JokerTrickster/common/error"
	"github.com/labstack/echo/v4"
)

// TokenChecker checks the JWT token in the request header
func TokenChecker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		accessToken := c.Request().Header.Get("tkn")
		if accessToken == "" {
			return _error.CreateError(ctx, string(_error.ErrBadParameter), _error.Trace(), "no access token in header", string(_error.ErrFromClient))
		}

		if err := VerifyToken(accessToken); err != nil {
			return err
		}

		uID, email, err := ParseToken(accessToken)
		if err != nil {
			return err
		}

		c.Set("uID", uID)
		c.Set("email", email)
		return next(c)
	}
}
