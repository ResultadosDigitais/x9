package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ResultadosDigitais/x9/log"
)

// ProcessRepositories is a work that receives a git repository and process it
// by clonning and analyzing.
func ProcessRepositories() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {

			for {
				repositoryID := <-session.Repositories
				repo, err := GetRepository(session, repositoryID)

				if err != nil {
					log.Error(fmt.Sprintf("Failed to retrieve repository %d: %s", repositoryID, err), nil)
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

// ProcessGists is a work that receives a gist and process it
// by clonning and analyzing.
func ProcessGists() {
	threadNum := *session.Options.Threads

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				gistURL := <-session.Gists
				processRepositoryOrGist(gistURL)
			}
		}(i)
	}
}

func processRepositoryOrGist(url string) {
	var (
		matches    []string
		matchedAny bool = false
	)

	dir := GetTempDir(GetHash(url))
	_, err := CloneRepository(session, url, dir)
	if err != nil {
		log.Debug(fmt.Sprintf("[%s] Cloning failed: %s", url, err.Error()), nil)
		os.RemoveAll(dir)
		return
	}
	log.Debug(fmt.Sprintf("[%s] Cloning in to %s", url, strings.Replace(dir, *session.Options.TempDirectory, "", -1)), nil)

	for _, file := range GetMatchingFiles(dir) {
		relativeFileName := strings.Replace(file.Path, *session.Options.TempDirectory, "", -1)
		relativeFileName = strings.SplitAfterN(relativeFileName, string(os.PathSeparator), 3)[2]
		if *session.Options.SearchQuery != "" {
			queryRegex := regexp.MustCompile(*session.Options.SearchQuery)
			for _, match := range queryRegex.FindAllSubmatch(file.Contents, -1) {
				hiddenMatch := obfuscate(string(match[0]))
				matches = append(matches, hiddenMatch)
			}
			if matches != nil {
				count := len(matches)
				m := strings.Join(matches, ", ")
				log.Result(session.Config.SlackWebhook, fmt.Sprintf(":warning: *Ooops I found something...*\n*Repository:* %s\n*Matches:* %d\n*Vulnerability:* %s\n*File:* %s\n*Values:* %s\n", url, count, "Search Query", relativeFileName, m))

				session.WriteToCsv([]string{url, "Search Query", relativeFileName, m})
			}
		} else {
			for _, signature := range session.Signatures {
				if matched, part := signature.Match(file); matched {
					matchedAny = true

					if part == PartContents {
						if matches = signature.GetContentsMatches(file); matches != nil {
							for i, _ := range matches {
								matches[i] = obfuscate(matches[i])
							}
							count := len(matches)
							m := strings.Join(matches, ", ")
							log.Result(session.Config.SlackWebhook, fmt.Sprintf(":warning: *Ooops I found something...*\n*Repository:* %s\n*Matches:* %d\n *Vulnerability:* %s\n*File:* %s\n*Values:* %s\n", url, count, signature.Name(), relativeFileName, m))

							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, m})
						}
					} else {
						if *session.Options.PathChecks {
							log.Result(session.Config.SlackWebhook, fmt.Sprintf(":warning: *Ooops I found something...*\n*Repository:* %s\n*File:* %s\n*Vulnerability:* %s\n", url, relativeFileName, signature.Name()))

							session.WriteToCsv([]string{url, signature.Name(), relativeFileName, ""})
						}

						if *session.Options.EntropyThreshold > 0 && file.CanCheckEntropy() {
							scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

							for scanner.Scan() {
								line := scanner.Text()

								if len(line) > 6 && len(line) < 100 {
									entropy := GetEntropy(scanner.Text())

									if entropy >= *session.Options.EntropyThreshold {
										log.Result(session.Config.SlackWebhook, fmt.Sprintf(":warning: *Ooops I found something...*\n*Repository:* %s\n*Vulnerability*: Potential secret in %s = %s", url, relativeFileName, scanner.Text()))

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

func obfuscate(text string) string {
	size := len(text)
	return text[0:2*size/3] + strings.Repeat("*", size/3)
}
