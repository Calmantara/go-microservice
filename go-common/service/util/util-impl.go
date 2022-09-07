package serviceutil

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UtilServiceImpl struct{}

func NewUtilService() UtilService {
	return &UtilServiceImpl{}
}

func (c *UtilServiceImpl) ContextBackground(ctx context.Context) (result context.Context) {
	result = context.Background()
	result = context.WithValue(result, logger.CorrelationKey.String(), ctx.Value(logger.CorrelationKey.String()))
	return result
}

func (c *UtilServiceImpl) SetCorrelationIdFromHeader(ctx *gin.Context) {
	corr := ctx.GetHeader(logger.CorrelationKey.String())
	if corr == "" {
		corr = uuid.New().String()
		val, exist := ctx.Get(logger.CorrelationKey.String())
		if exist {
			corr = val.(string)
		}
	}
	ctx.Set(logger.CorrelationKey.String(), corr)
}

func (c *UtilServiceImpl) UpsertCorrelationId(ctx context.Context, corrUid ...string) (ctxResult context.Context) {
	// check correlation id
	corr := (ctx).Value(logger.CorrelationKey.String())
	if corr == nil && len(corrUid) <= 0 {
		corr = uuid.New().String()
	} else if len(corrUid) > 0 {
		corr = corrUid[0]
	}
	(ctxResult) = context.WithValue(ctx, logger.CorrelationKey.String(), corr)
	return ctxResult
}

func (c *UtilServiceImpl) ObjectMapper(source interface{}, destination interface{}) (err error) {
	// encode
	byteObject, err := json.Marshal(&source)
	if err != nil {
		return err
	}
	// decode to object
	err = json.Unmarshal(byteObject, &destination)
	return err
}

func (c *UtilServiceImpl) InsertCorrelationIdToGrpc(ctx context.Context) (ctxResult context.Context) {
	corr := c.GetCorrelationIdFromContext(ctx)
	return metadata.AppendToOutgoingContext(ctx, logger.CorrelationKey.String(), corr)
}

func (c *UtilServiceImpl) GetCorrelationIdFromGrpc(ctx context.Context) (ctxResult context.Context) {
	meta, _ := metadata.FromIncomingContext(ctx)

	var corr string
	if len(meta[strings.ToLower(logger.CorrelationKey.String())]) > 0 {
		corr = meta[strings.ToLower(logger.CorrelationKey.String())][0]
	}
	return context.WithValue(ctx, logger.CorrelationKey.String(), corr)
}

func (c *UtilServiceImpl) GetCorrelationIdFromContext(ctx context.Context) (result string) {
	corr := ctx.Value(logger.CorrelationKey.String())
	// mappiug corr id
	c.ObjectMapper(corr, &result)
	if result == "" {
		result = uuid.New().String()
	}
	return result
}

// ### Response Util Function

func (c *UtilServiceImpl) UtilResponseBuilder(responseType model.ResponseType, data any) (response model.CommonResponseType) {
	switch responseType {
	case model.SUCCESS_OK_TYPE:
		response = model.CommonResponseType{
			HttpCode: http.StatusOK,
			CommonResponse: model.CommonResponse{
				Response: model.Response{
					ResponseCode:    "00",
					ResponseType:    model.SUCCESS_OK_TYPE,
					ResponseMessage: "request succeed"},
			},
		}

	case model.SUCCESS_ACCEPTED_TYPE:
		response = model.CommonResponseType{
			HttpCode: http.StatusAccepted,
			CommonResponse: model.CommonResponse{
				Response: model.Response{
					ResponseCode:    "01",
					ResponseType:    model.SUCCESS_ACCEPTED_TYPE,
					ResponseMessage: model.SUCCESS_ACCEPTED_MSG},
			},
		}
	}
	response.CommonResponse.ResponseData = data
	return response
}

func (c *UtilServiceImpl) UtilErrorResponseBuilder(errResponseType model.ResponseType, data any) (r model.CommonErrorResponseType) {
	switch errResponseType {
	case model.ERR_NO_WALLET_TYPE:
		r = model.CommonErrorResponseType{
			HttpCode: http.StatusBadRequest,
			CommonErrorResponse: model.CommonErrorResponse{
				Response: model.Response{
					ResponseCode:    "80",
					ResponseType:    model.ERR_NO_WALLET_TYPE,
					ResponseMessage: model.ERR_NO_WALLET_MSG},
			},
		}
	case model.ERR_BLOCKED_TYPE:
		r = model.CommonErrorResponseType{
			HttpCode: http.StatusBadRequest,
			CommonErrorResponse: model.CommonErrorResponse{
				Response: model.Response{
					ResponseCode:    "97",
					ResponseType:    model.ERR_BLOCKED_TYPE,
					ResponseMessage: model.ERR_BLOCKED_MSG},
			},
		}

	case model.ERR_BAD_REQUEST_TYPE:
		r = model.CommonErrorResponseType{
			HttpCode: http.StatusBadRequest,
			CommonErrorResponse: model.CommonErrorResponse{
				Response: model.Response{
					ResponseCode:    "98",
					ResponseType:    model.ERR_BAD_REQUEST_TYPE,
					ResponseMessage: model.ERR_BAD_REQUEST_MSG},
			},
		}

	default:
		r = model.CommonErrorResponseType{
			HttpCode: http.StatusInternalServerError,
			CommonErrorResponse: model.CommonErrorResponse{
				Response: model.Response{
					ResponseCode:    "99",
					ResponseType:    model.ERR_INTERNAL_TYPE,
					ResponseMessage: model.ERR_INTERNAL_MSG},
			},
		}
	}
	r.CommonErrorResponse.InvalidArgs = data
	return r
}

func (c *UtilServiceImpl) UtilErrorResponseSwitcher(ctx *gin.Context, errResponseType model.ResponseType, data any) {
	// go to error switcher
	r := c.UtilErrorResponseBuilder(errResponseType, data)
	ctx.AbortWithStatusJSON(r.HttpCode, r.CommonErrorResponse)
}

func (c *UtilServiceImpl) UtilResponseSwitcher(ctx *gin.Context, responseType model.ResponseType, data any) {
	// go to response switcher
	r := c.UtilResponseBuilder(responseType, data)
	ctx.JSONP(r.HttpCode, r.CommonResponse)
}
