package env

/*
	Context 생성 및 관리
*/

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type CtxValues struct {
	Method    string
	Url       string
	UserID    uint
	StartTime time.Time
	RequestID string
	Email     string
}

// CtxGenerate creates a context with custom values from Echo context
func CtxGenerate(c echo.Context) (context.Context, uint, string) {
	userID, _ := c.Get("uID").(uint)
	requestID, _ := c.Get("rID").(string)
	startTime, _ := c.Get("startTime").(time.Time)
	email, _ := c.Get("email").(string)
	req := c.Request()

	ctx := context.WithValue(req.Context(), "key", &CtxValues{
		Method:    req.Method,
		Url:       req.URL.Path,
		UserID:    userID,
		RequestID: requestID,
		StartTime: startTime,
		Email:     email,
	})
	return ctx, userID, email
}
