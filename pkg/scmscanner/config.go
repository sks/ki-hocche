package scmscanner

import (
	"net/http"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport/oauth2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sks/kihocche/pkg/trouble"
)

type SCMConfig struct {
	Token string
	Api   string
}

func (s SCMConfig) transport() http.RoundTripper {
	if s.Token == "" {
		return http.DefaultTransport
	}
	return &oauth2.Transport{
		Source: oauth2.StaticTokenSource(&scm.Token{
			Token: s.Token,
		}),
	}
}

func (s SCMConfig) client() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = s.transport()
	return retryClient.StandardClient()
}

type Configs []Config

func (c Configs) Scanners() ([]Scanner, error) {
	scanners := make([]Scanner, 0, len(c))
	for _, config := range c {
		scanner, err := config.Scanner()
		if err != nil {
			return nil, err
		}
		scanners = append(scanners, scanner)
	}
	return scanners, nil
}

type Config struct {
	Type   string    `yaml:"type"`
	Config SCMConfig `yaml:"config"`
}

type scannerConstructor func(SCMConfig) (Scanner, error)

type scannerFactory map[string]scannerConstructor

func (s scannerFactory) supported() []string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

func (c Config) Scanner() (Scanner, error) {
	supportedProviders := scannerFactory{
		"github": newGithubScanner,
	}
	scanner, ok := supportedProviders[c.Type]
	if !ok {
		return nil, trouble.New("unsupported provider", http.StatusNotImplemented, map[string]interface{}{
			"provider":  c.Type,
			"supported": supportedProviders.supported(),
		})
	}
	return scanner(c.Config)
}
