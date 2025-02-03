package formatter

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sks/kihocche/pkg/models"
	"github.com/sks/kihocche/pkg/trouble"
)

type FormatterType string

const (
	JSON FormatterType = "json"
	YAML FormatterType = "yaml"
	ICS  FormatterType = "ics"
)

type Config struct {
	Type   string `yaml:"type"`
	Output string
}

func (c Config) writer() (io.WriteCloser, error) {
	if c.Output == "" || c.Output == "stdout" {
		return os.Stdout, nil
	} else if c.Output == "stderr" {
		return os.Stderr, nil
	}
	writer, err := os.Create(c.Output)
	if err != nil {
		return nil, trouble.New("failed to open file", http.StatusInternalServerError, map[string]interface{}{
			"file": c.Output,
		})
	}
	return writer, nil
}

type Spec interface {
	Write(ctx context.Context, val models.Events) error
}

type spec interface {
	Write(writer io.Writer, val models.Events) error
}

func (c Config) Write(ctx context.Context, val models.Events) error {
	f, err := c.formatter()
	if err != nil {
		return err
	}
	writer, err := c.writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	return f.Write(writer, val)
}

func (c Config) formatter() (spec, error) {
	switch strings.ToLower(c.Type) {
	case string(JSON):
		return jsonFormatter{}, nil
	case string(YAML):
		return yamlFormatter{}, nil
	case string(ICS):
		return ICSFormatter{}, nil
	default:
		return nil, trouble.New("unsupported formatter", http.StatusNotImplemented, map[string]interface{}{
			"formatter": c.Type,
		})
	}
}
