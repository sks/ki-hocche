package trouble

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sks/kihocche/pkg/logger"
	"github.com/sks/kihocche/pkg/osutils"
	"github.com/sks/kihocche/pkg/tracer"
)

type trouble struct {
	ErrCode  string `json:"errCode,omitempty"`
	HttpCode int    `json:"httpCode,omitempty"`
	Extras   any    `json:"extras,omitempty"`
}

func New(errCode string, httpCode int, extras any) trouble {
	return trouble{
		ErrCode:  errCode,
		HttpCode: httpCode,
		Extras:   extras,
	}
}

func (e trouble) Error() string {
	return e.ErrCode
}

func (p trouble) JSON(w io.Writer) {
	_ = json.NewEncoder(w).Encode(p)
}

func WriteError(rw http.ResponseWriter, req *http.Request, err error) {
	problem, ok := toErr(err)
	if !ok {
		logger.GetLogger(req.Context()).Error("Error handling request", "error", err)
		problem = newMaskedError(tracer.TraceID(req.Context()))
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(problem.HttpCode)
	problem.JSON(rw)
}

func Convert(ctx context.Context, err error) trouble {
	problem, ok := toErr(err)
	if ok {
		return problem
	}
	logger.GetLogger(ctx).Error("an internal error occurred", "source", osutils.GetTraceInfo(3), "err", err)
	return newMaskedError(tracer.TraceID(ctx))
}

func newMaskedError(traceID string) trouble {
	return New("INTERNAL_SERVER_ERROR", http.StatusInternalServerError, map[string]any{
		"trace": traceID,
	})
}

func toErr(err error) (p trouble, ok bool) {
	ok = errors.As(err, &p)
	return p, ok
}

func IsTrouble(err error) bool {
	_, ok := toErr(err)
	return ok
}
