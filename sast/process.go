package sast

import (
	"fmt"
	"os"
	"time"

	"github.com/ResultadosDigitais/x9/git"
	"github.com/ResultadosDigitais/x9/log"
	"github.com/ResultadosDigitais/x9/util"
	"github.com/google/go-github/github"
)

type ProcessWorker struct {
	Session git.GithubSession
	Leaks   Leaks
	Events  chan *github.PullRequestEvent
}

func (pw *ProcessWorker) InitWorkers(w int) {
	for i := 0; i < w; i++ {
		go pw.ProcessEvent()
	}
}

func (pw *ProcessWorker) ProcessEvent() {
	for e := range pw.Events {
		repository, err := pw.Session.GetRepository(e.GetRepo().GetID())
		url := repository.GetCloneURL()
		if err != nil {
			log.Error(fmt.Sprintf("Error getting repository info: %s", url), map[string]interface{}{"error": err.Error()})
		}
		dir := util.GetTempDir(os.TempDir(), util.GetHash(url, time.Now().String()))
		if _, err := pw.Session.CloneRepository(url, dir); err != nil {
			log.Error(fmt.Sprintf("Error cloning repository: %s", url), map[string]interface{}{"error": err.Error()})

		}

		pw.Leaks.Test(url, dir)
		os.RemoveAll(dir)
	}
}
