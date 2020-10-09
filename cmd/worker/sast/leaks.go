package sast

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/ResultadosDigitais/x9/management"

	"github.com/ResultadosDigitais/x9/log"

	"github.com/ResultadosDigitais/x9/log/slack"
)

type Leaks struct {
	Config     Config
	Signatures []Signature
}

func (l *Leaks) Test(url, dir string) {
	log.Info(fmt.Sprintf("Testing repository %s", url), nil)
	for _, file := range GetMatchingFiles(dir, l.Config) {
		relativeFileName := strings.Replace(file.Path, os.TempDir(), "", -1)
		relativeFileName = strings.SplitAfterN(relativeFileName, string(os.PathSeparator), 3)[2]

		for _, signature := range l.Signatures {
			if matched, part := signature.Match(file); matched {

				if part == PartContents {
					l.processMatches(file, signature, url, relativeFileName)
				} else {
					fields := map[string]interface{}{
						"repo": url,
						"vuln": signature.Name(),
						"file": relativeFileName,
					}
					go insertVulnerability(url, signature.Name(), relativeFileName, "N/A")

					log.Info("Vulnerability found", fields)

					l.checkEntropy(file, url, relativeFileName)
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
	l.Signatures = GetSignatures(l.Config)

	return nil

}

func (l *Leaks) processMatches(file MatchFile, signature Signature, url, relativeFileName string) {
	if matches := signature.GetContentsMatches(file); matches != nil {
		for i, _ := range matches {
			matches[i] = obfuscate(matches[i])
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
		go insertVulnerability(url, signature.Name(), relativeFileName, m)
		log.Info("Vulnerability found", fields)
	}
}
func (l *Leaks) checkEntropy(file MatchFile, url, relativeFileName string) {
	if file.CanCheckEntropy(l.Config) {
		scanner := bufio.NewScanner(bytes.NewReader(file.Contents))

		for scanner.Scan() {
			line := scanner.Text()

			if len(line) > 6 && len(line) < 100 {
				entropy := getEntropy(scanner.Text())

				if entropy >= 5 {
					fields := map[string]interface{}{
						"repo":    url,
						"vuln":    "Potential secret",
						"file":    relativeFileName,
						"matches": 1,
						"values":  scanner.Text(),
					}
					go insertVulnerability(url, "Potential secret", relativeFileName, scanner.Text())
					log.Info("Vulnerability found", fields)
				}
			}
		}
	}
}

func getEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}

	for i := 0; i < 256; i++ {
		px := float64(strings.Count(data, string(byte(i)))) / float64(len(data))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}

	return entropy
}
func getVulnerabilityStructure(url, name, filename, values string) management.Vulnerability {
	return management.Vulnerability{
		Name:       name,
		Repository: url,
		FileName:   filename,
		Value:      values,
		Tool:       "X9 - Leak",
	}
}

// Obfuscate changes the last 1/3 string characters by *'s
func obfuscate(text string) string {
	size := len(text)
	return text[0:2*size/3] + strings.Repeat("*", len(text)-2*size/3)
}

func insertVulnerability(repository, vulnerability, fileName, values string) {
	vuln := getVulnerabilityStructure(repository, vulnerability, fileName, values)
	vuln, err := management.InsertVulnerability(vuln)
	if err != nil {
		log.Error("Database error", map[string]interface{}{"error": err.Error()})
	}
	if !vuln.FalsePositive {
		slack.Send(repository, vulnerability, fileName, values, vuln.ID, vuln.IssueURL)
	}
}
