//go:generate mockgen -source util.go -destination mock/util_mock.go -package mock

package serviceutil

import (
	"context"

	"github.com/Calmantara/go-common/model"
	"github.com/gin-gonic/gin"
)

type UtilService interface {
	// response builder
	UtilResponseBuilder(responseType model.ResponseType, data any) (response model.CommonResponseType)
	UtilErrorResponseBuilder(errResponseType model.ResponseType, data any) (response model.CommonErrorResponseType)
	UtilResponseSwitcher(ctx *gin.Context, responseType model.ResponseType, data any)
	UtilErrorResponseSwitcher(ctx *gin.Context, errResponseType model.ResponseType, data any)
	// correlation
	SetCorrelationIdFromHeader(ctx *gin.Context)
	UpsertCorrelationId(ctx context.Context, corrUid ...string) (ctxResult context.Context)
	GetCorrelationIdFromContext(ctx context.Context) (result string)
	// GRPC context
	InsertCorrelationIdToGrpc(ctx context.Context) (ctxResult context.Context)
	GetCorrelationIdFromGrpc(ctx context.Context) (ctxResult context.Context)
	// Mapping object can to struct
	ObjectMapper(source, destination any) (err error)
	// Generate new context background with correlation id
	ContextBackground(ctx context.Context) (result context.Context)
}
