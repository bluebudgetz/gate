package render

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	supportedContentTypes = []string{
		gin.MIMEYAML,
		"application/yaml",
		"text/yaml",
		gin.MIMEJSON,
		gin.MIMEXML,
		gin.MIMEXML2,
		gin.MIMEPlain,
		gin.MIMEHTML,
	}
)

func Render(c *gin.Context, statusCode int, data interface{}) {
	acceptedMimeType := c.NegotiateFormat(supportedContentTypes...)
	if acceptedMimeType == "" || len(c.Accepted) == 0 {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	switch acceptedMimeType {
	case gin.MIMEJSON:
		c.JSON(statusCode, data)

	case gin.MIMEHTML:
		c.Header("Content-Type", gin.MIMEHTML)
		c.Status(http.StatusOK) // send HTTP 200 since this is most probably a browser
		if _, err := c.Writer.WriteString("<!DOCTYPE html><html><body><pre>"); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		jsonEncoder := json.NewEncoder(c.Writer)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.SetIndent("", "  ")
		if err := jsonEncoder.Encode(data); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if _, err := c.Writer.WriteString("</pre></body></html>"); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case gin.MIMEXML, gin.MIMEXML2:
		c.XML(statusCode, data)

	case gin.MIMEPlain:
		c.Header("Content-Type", gin.MIMEPlain)
		c.Status(statusCode)
		jsonEncoder := json.NewEncoder(c.Writer)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.SetIndent("", "  ")
		if err := jsonEncoder.Encode(data); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case "text/yaml", "application/yaml", gin.MIMEYAML:
		c.YAML(statusCode, data)

	default:
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
}
