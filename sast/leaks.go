package sast

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ResultadosDigitais/x9/log"
	"github.com/ResultadosDigitais/x9/util"

	"github.com/ResultadosDigitais/x9/log/slack"
)

type Leaks struct {
	Signatures []Signature
	Config     Config
}

func (l *Leaks) Test(url, dir string) {
	log.Info(fmt.Sprintf("Testing repository %s", url), nil)
	var (
		matches []string
	)
	for _, file := range GetMatchingFiles(dir, l.Config) {
		relativeFileName := strings.Replace(file.Path, "/var", "", -1)
		relativeFileName = strings.SplitAfterN(relativeFileName, string(os.PathSeparator), 3)[2]

		for _, signature := range l.Signatures {
			if matched, part := signature.Match(file); matched {

				if part == PartContents {
					if matches = signature.GetContentsMatches(file); matches != nil {
						for i, _ := range matches {
							matches[i] = util.Obfuscate(matches[i])
						}
						count := len(matches)
						m := strings.Join(matches, ", ")

						fields := map[string]interface{}{
							"repo":    url,
							"matches": count,
							"vuln":    signature.Name(),
							"file":    relativeFileName,
							"values":  m,
						}
						slack.Send(fields)
						log.Info("Vulnerability found", fields)
					}
				} else {
					fields := map[string]interface{}{
						"repo": url,
						"vuln": signature.Name(),
						"file": relativeFileName,
					}
					slack.Send(fields)
					log.Info("Vulnerability found", fields)

					if file.CanCheckEntropy(l.Config) {
						scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

						for scanner.Scan() {
							line := scanner.Text()

							if len(line) > 6 && len(line) < 100 {
								entropy := util.GetEntropy(scanner.Text())

								if entropy >= 5 {
									fields := map[string]interface{}{
										"repo":    url,
										"vuln":    "Potential secret",
										"file":    relativeFileName,
										"matches": 1,
										"values":  scanner.Text(),
									}
									slack.Send(fields)
									log.Info("Vulnerability found", fields)

								}
							}
						}
					}
				}
			}
		}
	}
}

func (l *Leaks) GetLeaksConfig() error {

	config, err := ParseConfig()
	if err != nil {
		return err
	}
	l.Config = config
	return nil

}
