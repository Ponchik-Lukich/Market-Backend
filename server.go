package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"os"
	"strconv"
	"yandex-team.ru/bstask/api/routes"
	"yandex-team.ru/bstask/config"
)

type Server struct {
	DB   *sqlx.DB
	Echo *echo.Echo
}

func (s *Server) Initialize() {
	s.DB = config.ConnectDB()
	//err := config.DropTables(s.DB)
	//if err != nil {
	//	panic(err)
	//}
	err := config.InitializeTables(s.DB)
	if err != nil {
		panic(err)
	}
	s.Echo = s.SetupServer()
}

func (s *Server) SetupServer() *echo.Echo {
	e := echo.New()

	disableRateLimiterEnv := os.Getenv("DISABLE_RATE_LIMITER")
	disableRateLimiter, err := strconv.ParseBool(disableRateLimiterEnv)
	if err != nil {
		disableRateLimiter = false
	}

	if disableRateLimiter {
		routes.SetupRoutes(e, s.DB, 10000, 10000)
	} else {
		routes.SetupRoutes(e, s.DB, 10, 10)
	}

	return e
}

func main() {
	server := &Server{}
	server.Initialize()
	defer func() {
		err := server.DB.Close()
		if err != nil {
			panic(err)
		}
	}()
	server.Echo.Logger.Fatal(server.Echo.Start(":8080"))
}
