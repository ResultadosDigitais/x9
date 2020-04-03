package main

import (
	"github.com/ResultadosDigitais/x9/db"
	"github.com/ResultadosDigitais/x9/git"

	"github.com/ResultadosDigitais/x9/router"
	"github.com/ResultadosDigitais/x9/sast"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"

	"github.com/ResultadosDigitais/x9/config"
	"github.com/ResultadosDigitais/x9/log"
)

func main() {
	config.ParseConfig()

	log.Init()
	log.Info("X9 started...", nil)

	if err := db.GetDB(); err != nil {
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
	e.POST("/events", handler.Action)

	e.Logger.Fatal(e.Start(":3000"))

}
