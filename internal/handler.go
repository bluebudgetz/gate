package internal

import (
	"github.com/rs/zerolog/log"
	stdlog "log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

type (
	Handler struct {
		env string
	}
)

func NewHandler(env string) *Handler {
	return &Handler{env}
}

func (h *Handler) CreateReverseProxy(host string, port int) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			pathTokens := strings.SplitN(request.URL.Path[1:], "/", 3)
			version := pathTokens[0]
			var path string
			if len(pathTokens) >= 3 {
				path = pathTokens[2]
			}

			request.URL.Scheme = "http"
			request.URL.Host = host + ":" + strconv.Itoa(port)
			request.URL.Path = "/" + version + "/" + path
			request.Header.Set("User-Agent", "com.bluebudgetz.gate")

			if log.Logger.Debug().Enabled() {
				if bytes, err := httputil.DumpRequestOut(request, true); err != nil {
					log.Warn().Err(err).Msg("Failed dumping request")
				} else {
					log.Debug().Msg(string(bytes))
				}
			}
			// TODO: evaluate whether we need to set "X-Forwarded-For" and "X-Real-Ip"
			// request.Header.Set("X-Forwarded-For", r.RemoteAddr)
			// request.Header.Set("X-Real-Ip", r.RemoteAddr)
		},
		ErrorLog: stdlog.New(log.Logger, "", 0),
		ModifyResponse: func(response *http.Response) error {
			if log.Logger.Debug().Enabled() {
				if bytes, err := httputil.DumpResponse(response, true); err != nil {
					log.Warn().Err(err).Msg("Failed dumping response")
				} else {
					log.Debug().Msg(string(bytes))
				}
			}
			// TODO: evaluate whether we need to manually set the "Server" header
			response.Header.Del("Server")
			response.Header.Del("X-Content-Type-Options")
			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			http.Error(writer, "Internal error occurred.", http.StatusInternalServerError)
			log.Warn().Err(err).Msg("Error writing body from target back to client")
		},
	}
}
