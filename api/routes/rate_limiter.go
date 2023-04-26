package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"net/http"
	"os"
)

var disableRateLimiter = os.Getenv("DISABLE_RATE_LIMITER") == "true"

func newLimiter(r float64, b int) *rate.Limiter {
	if disableRateLimiter {
		return nil
	}
	return rate.NewLimiter(rate.Limit(r), b)
}

func withLimiter(limiter *rate.Limiter, handlerFunc func(echo.Context, *sqlx.DB) error, db *sqlx.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		if !disableRateLimiter && !limiter.Allow() {
			return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
		}
		return handlerFunc(c, db)
	}
}
