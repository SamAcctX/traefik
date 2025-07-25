package replacepath

import (
	"context"
	"net/http"
	"net/url"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/middlewares"
	"github.com/traefik/traefik/v3/pkg/middlewares/observability"
)

const (
	// ReplacedPathHeader is the default header to set the old path to.
	ReplacedPathHeader = "X-Replaced-Path"
	typeName           = "ReplacePath"
)

// ReplacePath is a middleware used to replace the path of a URL request.
type replacePath struct {
	next http.Handler
	path string
	name string
}

// New creates a new replace path middleware.
func New(ctx context.Context, next http.Handler, config dynamic.ReplacePath, name string) (http.Handler, error) {
	middlewares.GetLogger(ctx, name, typeName).Debug().Msg("Creating middleware")

	return &replacePath{
		next: next,
		path: config.Path,
		name: name,
	}, nil
}

func (r *replacePath) GetTracingInformation() (string, string) {
	return r.name, typeName
}

func (r *replacePath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.RawPath
	if currentPath == "" {
		currentPath = req.URL.EscapedPath()
	}

	req.Header.Add(ReplacedPathHeader, currentPath)
	req.URL.RawPath = r.path

	var err error
	req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
	if err != nil {
		middlewares.GetLogger(context.Background(), r.name, typeName).Error().Msgf("Unable to unescape url raw path %q: %v", req.URL.RawPath, err)
		observability.SetStatusErrorf(req.Context(), "Unable to unescape url raw path %q: %v", req.URL.RawPath, err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.RequestURI = req.URL.RequestURI()

	r.next.ServeHTTP(rw, req)
}
