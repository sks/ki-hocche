package caldav

import (
	"net/http"

	"github.com/sks/kihocche/pkg/auth"
	"github.com/sks/kihocche/pkg/trouble"
)

type Server struct {
	auth auth.AuthService
}

func (s Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, err := s.auth.CurrentUserPrincipal(req.Context())
	if err != nil {
		trouble.WriteError(rw, req, err)
		return
	}
	trouble.WriteError(rw, req, trouble.New("not implemented", http.StatusNotImplemented, nil))
}
