package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sks/kihocche/pkg/models"
	"gopkg.in/yaml.v3"
)

type jsonFormatter struct {
}

func (j jsonFormatter) Write(writer io.Writer, val models.Events) error {
	err := json.NewEncoder(writer).Encode(val)
	if err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}
	return nil
}

// yamlFormatter is a yaml formatter
type yamlFormatter struct {
}

// Write writes the value to the output
func (y yamlFormatter) Write(writer io.Writer, val models.Events) error {
	err := yaml.NewEncoder(writer).Encode(val)
	if err != nil {
		return fmt.Errorf("failed to write yaml: %w", err)
	}
	return nil
}
