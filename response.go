package httputil

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-courier/statuserror"
	"net/http"
)

type RESP struct {
	ServiceCode int
	Content     any
	Err         error
}

func (resp *RESP) Output(ctx *gin.Context) {
	if resp.Err == nil {
		ctx.JSON(http.StatusOK, resp.Content)
		return
	}

	var e statuserror.StatusError
	if ok := errors.As(resp.Err, &e); ok {
		ctx.AbortWithStatusJSON(e.StatusErr().StatusCode(), ErrorRESP{
			Code: e.StatusErr().Code,
			Msg:  e.StatusErr().Msg,
		})
		return

	}

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRESP{
		Code: http.StatusInternalServerError*1e6 + resp.ServiceCode,
		Msg:  resp.Err.Error(),
	})
	return
}

type ErrorRESP struct {
	Code int    `json:"type"`
	Msg  string `json:"msg"`
}
