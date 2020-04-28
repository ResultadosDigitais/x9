package main

import (
	"github.com/ResultadosDigitais/x9/cmd/worker/git"
	"github.com/ResultadosDigitais/x9/db"

	"github.com/ResultadosDigitais/x9/cmd/worker/router"
	"github.com/ResultadosDigitais/x9/cmd/worker/sast"
	"github.com/google/go-github/github"
	"github.com/labstack/echo/v4"

	"github.com/ResultadosDigitais/x9/config"
	"github.com/ResultadosDigitais/x9/log"
)

func main() {
	log.Init()
	log.Info("X9 started...", nil)

	if err := config.ParseConfig(); err != nil {
		log.Fatal("Error on parsing config", map[string]interface{}{"error": err.Error()})
	}

	if err := db.GetDB(config.Opts.DatabaseConfig); err != nil {
		log.Fatal("Database connection error", map[string]interface{}{"error": err.Error()})
	}
	db.InitTables()

	githubSession := git.GithubSession{}
	if err := githubSession.InitClient(); err != nil {
		log.Fatal(err.Error(), nil)
	}

	leaks := sast.Leaks{}
	if err := leaks.GetLeaksConfig(); err != nil {
		log.Fatal(err.Error(), nil)

	}

	eventsChannel := make(chan *github.PullRequestEvent)
	processWorker := sast.ProcessWorker{
		Session: githubSession,
		Leaks:   leaks,
		Events:  eventsChannel,
	}
	processWorker.InitWorkers(*config.Opts.Threads)

	handler := router.Handler{
		Process: eventsChannel,
		Session: &githubSession,
	}
	e := echo.New()

	e.GET("/healthcheck", handler.HealthCheck)
	e.POST("/events", handler.Event)
	e.POST("/action", handler.Action)

	e.Logger.Fatal(e.Start(":3000"))

}
