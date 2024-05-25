package response

import (
	"github.com/gin-gonic/gin"
	"ktserver/internal/pkg/bizerr"
	"net/http"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func wrapBaseResponse(v any) Response[any] {
	var resp Response[any]
	switch data := v.(type) {
	case *bizerr.CodeMsg:
		resp.Code = data.Code
		resp.Msg = data.Msg
		break
	case bizerr.CodeMsg:
		resp.Code = data.Code
		resp.Msg = data.Msg
	//case *status.Status: // grpc status
	//	resp.Code = int(data.Code())
	//	resp.Msg = data.Message()
	case error:
		resp.Code = bizerr.UnknownError.Code
		resp.Msg = data.Error()
	default:
		resp.Code = bizerr.Success.Code
		resp.Msg = bizerr.Success.Msg
		resp.Data = v
	}
	return resp
}

func Result(c *gin.Context, v any) {
	c.JSON(http.StatusOK, wrapBaseResponse(v))
}

func ResultWithStatus(c *gin.Context, status int, v any) {
	c.JSON(status, wrapBaseResponse(v))
}
