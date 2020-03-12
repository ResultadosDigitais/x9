package main

import (
	"fmt"

	"github.com/ResultadosDigitais/x9/git"
	"github.com/ResultadosDigitais/x9/router"
	"github.com/ResultadosDigitais/x9/sast"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"

	"github.com/ResultadosDigitais/x9/log"
)

func main() {

	log.Init()
	log.Info("X9 started...", nil)

	eventsChannel := make(chan *github.PullRequestEvent)

	githubSession := git.GithubSession{}
	if err := githubSession.InitClient(); err != nil {
		log.Fatal(err.Error(), nil)
	}
	fmt.Println(githubSession)

	leaks := sast.Leaks{}
	if err := leaks.GetLeaksConfig(); err != nil {
		log.Fatal(err.Error(), nil)

	}

	processWorker := sast.ProcessWorker{
		Session: githubSession,
		Leaks:   leaks,
		Events:  eventsChannel,
	}

	processWorker.InitWorkers(5)

	handler := router.Handler{
		Process: eventsChannel,
	}

	e := echo.New()

	e.GET("/healthcheck", handler.HealthCheck)
	e.POST("/events", handler.Event)

	e.Logger.Fatal(e.Start(":3000"))

}
