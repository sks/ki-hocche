package encodeutils

import (
	"context"
	"encoding/json"

	"github.com/sks/kihocche/pkg/logger"
)

func MustToJSON(ctx context.Context, v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		logger.GetLogger(ctx).Error("Error marshalling to JSON", "v", v, "error", err)
	}
	return data
}
