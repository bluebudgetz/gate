package infra

import (
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewRouter(cfg config.Config) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false // adding/removing trailing slash automatically is confusing; disabling
	router.RedirectFixedPath = true      // cleans up ugly paths like "../" and "//"
	router.HandleMethodNotAllowed = true // detects & handles MethodNotAllowed situations (HTTP/405) automatically
	router.ForwardedByClientIP = true    // infers IP from "X-Forwarded-For" & "X-Real-Ip" headers (if any)
	router.AppEngine = false             // disables Google AppEngine integration
	router.MaxMultipartMemory = 32 << 20 // 32 MB (TODO: externalize max-multipart-memory)
	router.Use(
		// Request ID & Server headers
		func(c *gin.Context) {
			id := c.GetHeader("X-Request-ID")
			if id == "" {
				id = uuid.New().String()
			}
			c.Set("X-Request-ID", id)
			c.Writer.Header().Set("X-Request-ID", id)
			c.Writer.Header().Set("Server", "Bluebudgetz/gate")
		},

		// Request log
		ginLogger(cfg.HTTP.DisableLogRequests),

		// Error handlers
		ginHandlerForErrorsOfType(http.StatusBadRequest, gin.ErrorTypeBind),
		ginHandlerForLastErrorOfType(
			http.StatusInternalServerError,
			gin.ErrorTypePrivate,
			"Oops, that's embarrassing! Something has gone wrong, we're on it.",
		),

		// Recover from panics
		ginRecover(),

		// CORS
		cors.New(cors.Config{
			AllowOrigins:     cfg.HTTP.CORS.AllowOrigins,
			AllowMethods:     cfg.HTTP.CORS.AllowMethods,
			AllowHeaders:     cfg.HTTP.CORS.AllowHeaders,
			ExposeHeaders:    cfg.HTTP.CORS.ExposeHeaders,
			AllowCredentials: cfg.HTTP.CORS.AllowCredentials,
			MaxAge:           time.Duration(cfg.HTTP.CORS.MaxAge) * time.Second,
		}),

		// Secure
		secure.New(secure.Config{
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			IENoOpen:              true,
			ReferrerPolicy:        "strict-origin-when-cross-origin",
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		}),

		// GZip support
		gzip.Gzip(cfg.HTTP.GZipLevel),
	)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, "OK\n")
		c.Writer.Flush()
	})

	// TODO: add JWT authentication

	return router, nil
}
