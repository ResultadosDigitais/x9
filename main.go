package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ResultadosDigitais/policeman/core"
	"github.com/ResultadosDigitais/policeman/log"
)

var session = core.GetSession()

func ProcessRepositories() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {

			for {
				repositoryId := <-session.Repositories
				repo, err := core.GetRepository(session, repositoryId)

				if err != nil {
					log.Error(fmt.Sprintf("Failed to retrieve repository %d: %s", repositoryId, err), nil)
					continue
				}

				if repo.GetPermissions()["pull"] &&
					uint(repo.GetStargazersCount()) >= *session.Options.MinimumStars &&
					uint(repo.GetSize()) < *session.Options.MaximumRepositorySize {

					processRepositoryOrGist(repo.GetCloneURL())
				}
			}
		}(i)
	}
}

func ProcessGists() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				gistUrl := <-session.Gists
				processRepositoryOrGist(gistUrl)
			}
		}(i)
	}
}

func processRepositoryOrGist(url string) {
	var (
		matches    []string
		matchedAny bool = false
	)

	dir := core.GetTempDir(core.GetHash(url))
	_, err := core.CloneRepository(session, url, dir)
	if err != nil {
		log.Debug(fmt.Sprintf("[%s] Cloning failed: %s", url, err.Error()), nil)
		os.RemoveAll(dir)
		return
	}
	log.Debug(fmt.Sprintf("[%s] Cloning in to %s", url, strings.Replace(dir, *session.Options.TempDirectory, "", -1)), nil)

	for _, file := range core.GetMatchingFiles(dir) {
		relativeFileName := strings.Replace(file.Path, *session.Options.TempDirectory, "", -1)

		if *session.Options.SearchQuery != "" {
			queryRegex := regexp.MustCompile(*session.Options.SearchQuery)
			for _, match := range queryRegex.FindAllSubmatch(file.Contents, -1) {
				matches = append(matches, string(match[0]))
			}

			if matches != nil {
				count := len(matches)
				m := strings.Join(matches, ", ")
				log.Result(session.Config.SlackWebhook, fmt.Sprintf("[%s] %d %s for %s in file %s: %s", url, count, core.Pluralize(count, "match", "matches"), "Search Query", relativeFileName, m))

				session.WriteToCsv([]string{url, "Search Query", relativeFileName, m})
			}
		} else {
			for _, signature := range session.Signatures {
				if matched, part := signature.Match(file); matched {
					matchedAny = true

					if part == core.PartContents {
						if matches = signature.GetContentsMatches(file); matches != nil {
							count := len(matches)
							m := strings.Join(matches, ", ")
							log.Result(session.Config.SlackWebhook, fmt.Sprintf("[%s] %d %s for %s in file %s: %s", url, count, core.Pluralize(count, "match", "matches"), signature.Name(), relativeFileName, m))

							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, m})
						}
					} else {
						if *session.Options.PathChecks {
							log.Result(session.Config.SlackWebhook, fmt.Sprintf("[%s] Matching file %s for %s", url, relativeFileName, signature.Name()))

							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, ""})
						}

						if *session.Options.EntropyThreshold > 0 && file.CanCheckEntropy() {
							scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

							for scanner.Scan() {
								line := scanner.Text()

								if len(line) > 6 && len(line) < 100 {
									entropy := core.GetEntropy(scanner.Text())

									if entropy >= *session.Options.EntropyThreshold {
										log.Info(fmt.Sprintf("[%s] Potential secret in %s = %s", url, relativeFileName, scanner.Text()), nil)

										session.WriteToCsv([]string{url, signature.Name(), relativeFileName, scanner.Text()})
									}
								}
							}
						}
					}
				}
			}
		}

		if !matchedAny {
			os.Remove(file.Path)
		}
	}

	if !matchedAny {
		os.RemoveAll(dir)
	}
}

func main() {
	log.Init()
	log.Info(fmt.Sprintf("%s v%s started. Loaded %d signatures. Using %d GitHub tokens and %d threads. Work dir: %s", core.Name, core.Version, len(session.Signatures), len(session.Clients), *session.Options.Threads, *session.Options.TempDirectory), nil)

	if *session.Options.SearchQuery != "" {
		log.Info(fmt.Sprintf("Search Query '%s' given. Only returning matching results.", *session.Options.SearchQuery), nil)

	}

	go core.GetRepositories(session)
	go ProcessRepositories()

	if *session.Options.ProcessGists {
		go core.GetGists(session)
		go ProcessGists()
	}

	log.Info("Press Crt+C to stop and exit. \n", nil)
	select {}
}
