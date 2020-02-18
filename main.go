package main

import (
	"fmt"

	"github.com/ResultadosDigitais/x9/core"
	"github.com/ResultadosDigitais/x9/log"
)

func main() {
	session := core.GetSession()
	log.Init()
	log.Info(fmt.Sprintf("%s v%s started. Loaded %d signatures. Using %d GitHub tokens and %d threads. Work dir: %s", core.Name, core.Version, len(session.Signatures), len(session.Clients), *session.Options.Threads, *session.Options.TempDirectory), nil)

	if *session.Options.SearchQuery != "" {
		log.Info(fmt.Sprintf("Search Query '%s' given. Only returning matching results.", *session.Options.SearchQuery), nil)

	}

	go core.GetRepositories(session)
	go core.ProcessRepositories()

	if *session.Options.ProcessGists {
		go core.GetGists(session)
		go core.ProcessGists()
	}

	log.Info("Press Crt+C to stop and exit. \n", nil)
	select {}
}
