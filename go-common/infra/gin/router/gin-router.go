//go:generate mockgen -source gin-router.go -destination mock/gin-router_mock.go -package mock

package ginrouter

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	group "github.com/Calmantara/go-common/infra/gin/group"
	"github.com/Calmantara/go-common/logger"
)

type GinRouter interface {
	GROUP(groupPath string, handlers ...gin.HandlerFunc) group.GinGroup
	USE(middleware ...gin.HandlerFunc)
	GET(path string, handler ...gin.HandlerFunc)
	POST(path string, handler ...gin.HandlerFunc)
	PUT(path string, handler ...gin.HandlerFunc)
	DELETE(path string, handler ...gin.HandlerFunc)
	SERVE(ops ...Option)
}

type GinRouterImpl struct {
	sugar *zap.SugaredLogger
	Gin   *gin.Engine
	Port  string
}

type Option func(g *GinRouterImpl)

func NewGinRouter(sugar logger.CustomLogger) GinRouter {
	sugar.Logger().Info("Initialization gin router")
	return &GinRouterImpl{
		Gin: gin.New()}
}

func (gr *GinRouterImpl) GROUP(groupPath string, handlers ...gin.HandlerFunc) group.GinGroup {
	return group.NewGinGroup(gr.Gin.Group(groupPath, handlers...))
}

func (gr *GinRouterImpl) USE(middleware ...gin.HandlerFunc) {
	gr.Gin.Use(middleware...)
}

func (gr *GinRouterImpl) GET(path string, handler ...gin.HandlerFunc) {
	gr.Gin.GET(path, handler...)
}

func (gr *GinRouterImpl) POST(path string, handler ...gin.HandlerFunc) {
	gr.Gin.POST(path, handler...)
}

func (gr *GinRouterImpl) PUT(path string, handler ...gin.HandlerFunc) {
	gr.Gin.PUT(path, handler...)
}

func (gr *GinRouterImpl) DELETE(path string, handler ...gin.HandlerFunc) {
	gr.Gin.DELETE(path, handler...)
}

func (gr *GinRouterImpl) SERVE(ops ...Option) {

	// iterate all option method
	for _, v := range ops {
		v(gr)
	}

	fmt.Printf("Running application in %v", gr.Port)
	gr.Gin.Run(":" + gr.Port)
}

func WithPort(port string) Option {
	return func(g *GinRouterImpl) {
		g.Port = port
	}
}
