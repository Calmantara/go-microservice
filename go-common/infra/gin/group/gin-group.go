//go:generate mockgen -source gin-group.go -destination mock/gin-group_mock.go -package mock

package gingroup

import "github.com/gin-gonic/gin"

type GinGroup interface {
	GET(path string, handler ...gin.HandlerFunc)
	POST(path string, handler ...gin.HandlerFunc)
	PUT(path string, handler ...gin.HandlerFunc)
	DELETE(path string, handler ...gin.HandlerFunc)
}

type GinGroupImpl struct {
	Rg *gin.RouterGroup
}

func NewGinGroup(rg *gin.RouterGroup) GinGroup {
	return &GinGroupImpl{
		Rg: rg,
	}
}

func (gr *GinGroupImpl) GET(path string, handler ...gin.HandlerFunc) {
	gr.Rg.GET(path, handler...)
}

func (gr *GinGroupImpl) POST(path string, handler ...gin.HandlerFunc) {
	gr.Rg.POST(path, handler...)
}

func (gr *GinGroupImpl) PUT(path string, handler ...gin.HandlerFunc) {
	gr.Rg.PUT(path, handler...)
}

func (gr *GinGroupImpl) DELETE(path string, handler ...gin.HandlerFunc) {
	gr.Rg.DELETE(path, handler...)
}
