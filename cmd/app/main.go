package main

import (
	"github.com/ResultadosDigitais/x9/db"
	"github.com/gorilla/sessions"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/ResultadosDigitais/x9/cmd/app/auth"
	"github.com/ResultadosDigitais/x9/cmd/app/handler"
	"github.com/ResultadosDigitais/x9/config"
	"github.com/ResultadosDigitais/x9/log"
)

func main() {

	log.Init()
	log.Info("X9 app started...", nil)

	if err := config.ParseAppConfig(); err != nil {
		log.Fatal("Error on parsing config", map[string]interface{}{"error": err.Error()})
	}

	if err := db.GetDB(config.AppOpts.DatabaseConfig); err != nil {
		log.Fatal("Database connection error", map[string]interface{}{"error": err.Error()})
	}
	if err := auth.InitOIDC(); err != nil {
		log.Fatal("OIDC configuration failed", map[string]interface{}{"error": err.Error()})

	}
	t := handler.InitTemplate()
	h := handler.InitHandler()
	e := echo.New()
	e.Renderer = &t
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.AppOpts.ApplicationSecretKey))))
	e.Static("/static", "templates/css")
	e.GET("/login", h.Login)
	e.GET("/auth/google/callback", h.OIDCCallback)
	e.GET("/auth", h.OIDCAuth)

	e.Logger.Fatal(e.Start(":3000"))

}
