//go:build disable_ratelimiter
// +build disable_ratelimiter

// routes/rate_limiter_disabled.go

package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

func newLimiter(r float64, d int) *rate.Limiter {
	return nil
}

func withLimiter(limiter *rate.Limiter, handlerFunc func(echo.Context, *sqlx.DB) error, db *sqlx.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		return handlerFunc(c, db)
	}
}
