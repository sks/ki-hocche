package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/sks/kihocche/pkg/logger"
	"github.com/sks/kihocche/pkg/trouble"
	"gocloud.dev/secrets"
)

type Router struct {
	secrets *secrets.Keeper
}

func NewRouter(secrets *secrets.Keeper) *Router {
	return &Router{
		secrets: secrets,
	}
}

// GenerateLinkPage renders a page with a form to generate a link
func (Router) GenerateLinkPage(w http.ResponseWriter, r *http.Request) {
	// html template with a form that asks for the scm type, scm url, and the filter
	// the form should have a submit button
	// the form should be submitted to /generate_link
	template := `
	<!DOCTYPE html>
	<html>
		<head>
			<title>Generate Link</title>
		</head>
		<body>
			<form action="/generate_link" method="POST">
				<h1>Generate Link</h1>
				<div>
					<label for="type">SCM Type</label>
					<select id="configs" name="type">
						<option value="github">GitHub</option>
					</select>
				</div>
				<div>
					<label for="secret">Secret</label>
					<input type="password" id="token" name="token" value="" placeholder="Token">
				</div>
				<div>
					<label for="Since">Since</label>
					<input type="text" id="since" name="since" value="7d">
				</div>
				<div>
					<label for="repos">Repos [comma seperated]</label>
					<input type="text" id="repos" name="repos" value="">
				</div>
				<div>
					<label for="namespace">Namespace [comma seperated]</label>
					<input type="text" id="namespace" name="namespace" value="">
				</div>
				<div>
					<label for="events">Events [comma seperated]</label>
					<input type="text" id="events" name="events" value="push,pull_request,release">
				</div>
				<input type="submit" value="Generate Link">
			</form>
		</body>
	</html>`
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(template))
}

func (router Router) GenerateLink(w http.ResponseWriter, r *http.Request) {
	var request webRequest
	if r.Method == http.MethodGet {
		request, _ = router.parseQuery(r)
	} else {
		request = router.parseForm(r)
	}
	key, err := router.encrypt(r.Context(), request)
	if err != nil {
		// write error
		trouble.WriteError(w, r, err)
		return
	}
	// write the key
	// this URL Should be /generate_link, make it /subscribe.ics?key=<key> and show it in the response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("Your link is: /subscribe.ics?key=%s", key)))
}

func (r Router) encrypt(ctx context.Context, request webRequest) (string, error) {
	// encrypt the scmConfig and filter
	// return the encrypted string
	result, err := r.secrets.Encrypt(ctx, request.byte())
	if err != nil {
		return "", fmt.Errorf("failed to encrypt the key: %w", err)
	}
	// base64 encode the result
	return base64.URLEncoding.EncodeToString(result), nil
}

func (r Router) decrypt(ctx context.Context, encrypted string) (webRequest, error) {
	// base64 decode the encrypted string
	// decrypt the result
	// return the webRequest
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		logger.GetLogger(ctx).Error("failed to decode the key", "key", encrypted, "error", err)
		return webRequest{}, fmt.Errorf("failed to decode the key: %w", err)
	}
	result, err := r.secrets.Decrypt(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("failed to decrypt the key", "error", err)
		return webRequest{}, fmt.Errorf("failed to decrypt the key: %w", err)
	}
	return newWebReq(result), nil
}
