package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"net/http"
)

func newLimiter(r float64, s int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(r), s)
}

func withLimiter(limiter *rate.Limiter, handlerFunc func(echo.Context, *sqlx.DB) error, db *sqlx.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		if !limiter.Allow() {
			return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
		}
		return handlerFunc(c, db)
	}
}
