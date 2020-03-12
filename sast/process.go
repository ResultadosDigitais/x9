package sast

import (
	"os"
	"time"

	"github.com/ResultadosDigitais/x9/git"
	"github.com/ResultadosDigitais/x9/util"
	"github.com/google/go-github/github"
)

type ProcessWorker struct {
	Session git.GithubSession
	Leaks   Leaks
	Events  chan *github.PullRequestEvent
}

func (pw *ProcessWorker) InitWorkers(w int) {
	for w := 1; w <= 3; w++ {
		go pw.ProcessEvent()
	}
}

func (pw *ProcessWorker) ProcessEvent() {
	for e := range pw.Events {
		repository, err := pw.Session.GetRepository(e.GetRepo().GetID())
		url := repository.GetURL()
		if err != nil {

		}
		dir := util.GetTempDir("/tmp", util.GetHash(url, time.Now().String()))
		_, err = pw.Session.CloneRepository(url, dir)

		pw.Leaks.Test(url, dir)
		os.RemoveAll(dir)
	}
}
