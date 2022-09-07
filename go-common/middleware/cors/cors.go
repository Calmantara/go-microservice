package cors

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type CorsMiddleware interface {
	Cors(ctx *gin.Context)
}

type CorsMiddlewareImpl struct {
	Origin           []string
	AdditionalHeader []string
}

type Option func(*CorsMiddlewareImpl)

func NewCorsMiddleware(corsOps ...Option) CorsMiddleware {
	corsMiddleware := CorsMiddlewareImpl{
		Origin:           []string{"*"},
		AdditionalHeader: defaultAcceptHeader(),
	}
	// loop over cors options function
	for _, val := range corsOps {
		val(&corsMiddleware)
	}

	return &corsMiddleware
}

func (x *CorsMiddlewareImpl) Cors(c *gin.Context) {
	c.Writer.Header().Set(`Access-Control-Allow-Origin`, strings.Join(x.Origin, ","))
	c.Writer.Header().Set(`Access-Control-Allow-Credentials`, `true`)
	c.Writer.Header().Set(`Access-Control-Allow-Headers`, strings.Join(x.AdditionalHeader, ","))
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func WithAdditionalOrigin(origins []string) Option {
	return func(cmi *CorsMiddlewareImpl) {
		cmi.Origin = origins
	}
}

func WithAdditionalHeader(Header []string) Option {
	return func(cmi *CorsMiddlewareImpl) {
		cmi.AdditionalHeader = append(cmi.AdditionalHeader, Header...)
	}
}

func defaultAcceptHeader() []string {
	return []string{`Content-Type`,
		`Content-Length`,
		`Accept-Encoding`,
		`X-CSRF-Token`,
		`Authorization`,
		`accept`,
		`origin`,
		`Cache-Control`,
		`X-Requested-With`,
		`Static-Token`,
		`Correlation-ID`}
}
