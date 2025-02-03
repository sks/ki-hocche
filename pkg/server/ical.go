package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sks/kihocche/pkg/formatter"
	"github.com/sks/kihocche/pkg/logger"
	"github.com/sks/kihocche/pkg/models"
	"github.com/sks/kihocche/pkg/scmscanner"
	"github.com/sks/kihocche/pkg/trouble"
	"golang.org/x/sync/errgroup"
)

type webRequest struct {
	Configs scmscanner.Configs `json:"configs"`
	Filter  scmscanner.Filter  `json:"filter"`
}

func (w webRequest) byte() []byte {
	data, _ := json.Marshal(w)
	return data
}

func newWebReq(data []byte) webRequest {
	var w webRequest
	_ = json.Unmarshal(data, &w)
	return w
}

func (router *Router) ServeICal(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger(r.Context()).With("handler", "ical")
	logger.Info("serving ical")
	defer func(start time.Time) {
		logger.Info("request served", "duration", time.Since(start).String())
	}(time.Now())
	webRequest, err := router.parseQuery(r)
	if err != nil {
		trouble.WriteError(w, r, err)
		return
	}
	scanners, err := webRequest.Configs.Scanners()
	if err != nil {
		trouble.WriteError(w, r, err)
		return
	}
	result := make(models.Events, 0)
	errGroup, ctx := errgroup.WithContext(r.Context())
	for _, scanner := range scanners {
		scanner := scanner
		errGroup.Go(func() error {
			events, err := scanner.Get(ctx, webRequest.Filter)
			if err != nil {
				return err
			}
			result = append(result, events...)
			return nil
		})
	}
	err = errGroup.Wait()
	if err != nil {
		trouble.WriteError(w, r, err)
		return
	}
	w.Header().Add("Content-Type", "text/calendar")
	w.WriteHeader(http.StatusOK)
	err = formatter.ICSFormatter{}.Write(w, result)
	if err != nil {
		trouble.WriteError(w, r, err)
		return
	}
}

func (router *Router) parseForm(r *http.Request) webRequest {
	_ = r.ParseForm()
	scmConfig := scmscanner.Config{
		Type: r.Form.Get("type"),
		Config: scmscanner.SCMConfig{
			Api:   r.Form.Get("api"),
			Token: r.Form.Get("token"),
		},
	}
	since, err := time.ParseDuration(r.Form.Get("since"))
	if err != nil {
		since = 7 * 24 * time.Hour
	}
	filter := scmscanner.Filter{
		Since:  since,
		Repos:  strings.Split(r.Form.Get("repos"), ","),
		Events: strings.Split(r.Form.Get("events"), ","),
	}
	return webRequest{
		Configs: scmscanner.Configs{scmConfig},
		Filter:  filter,
	}
}

func (router *Router) parseQuery(r *http.Request) (webRequest, error) {
	query := r.URL.Query()
	if key := query.Get("key"); key != "" {
		return router.decrypt(r.Context(), key)
	}
	scmConfig := scmscanner.Config{
		Type: query.Get("type"),
		Config: scmscanner.SCMConfig{
			Api:   query.Get("api"),
			Token: query.Get("token"),
		},
	}
	since, err := time.ParseDuration(query.Get("since"))
	if err != nil {
		since = 7 * 24 * time.Hour
	}
	filter := scmscanner.Filter{
		Since:  since,
		Repos:  query["repos"],
		Events: query["events"],
	}
	return webRequest{
		Configs: scmscanner.Configs{scmConfig},
		Filter:  filter,
	}, nil
}
