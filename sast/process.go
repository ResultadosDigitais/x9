package sast

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ResultadosDigitais/x9/config"
	"github.com/ResultadosDigitais/x9/git"
	"github.com/ResultadosDigitais/x9/log"
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
		if isBlacklisted(url, config.GetBlacklistedRepositories()) {
			continue
		}
		if err != nil {
			log.Error(fmt.Sprintf("Error getting repository info: %s", url), map[string]interface{}{"error": err.Error()})
		}
		dir := getDir(url)
		if _, err := pw.Session.CloneRepository(url, dir); err != nil {
			log.Error(fmt.Sprintf("Error cloning repository: %s", url), map[string]interface{}{"error": err.Error()})

		}

		pw.Leaks.Test(url, dir)
		os.RemoveAll(dir)
	}
}

func isBlacklisted(repo string, blackListedRepositories []string) bool {
	for _, blackListedRepo := range blackListedRepositories {
		if match, err := regexp.MatchString(`^.*(/`+blackListedRepo+`/).*$`, repo); err != nil {
			log.Error("Regexp error compile", map[string]interface{}{"error": err.Error()})
		} else if match == true {
			return true
		}
	}
	return false

}

func getDir(url string) string {
	folderName := getHash(url, time.Now().String())
	return getTempDir(os.TempDir(), folderName)

}

func getTempDir(tempDir, suffix string) string {
	dir := filepath.Join(tempDir, suffix)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	} else {
		os.RemoveAll(dir)
	}

	return dir
}

// getHash receives one or more strings and returns
// a sha1 hash of the concatenation of all of them
func getHash(values ...string) string {
	str := values[0]
	for i := 1; i < len(values); i++ {
		str = fmt.Sprintf("%s%s", str, values[i])
	}
	hashString := sha1.New()
	hashString.Write([]byte(str))
	return hex.EncodeToString(hashString.Sum(nil))
}
