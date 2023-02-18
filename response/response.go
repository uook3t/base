package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewGinResponse(code int32, data interface{}) *GinResponse {
	return &GinResponse{
		Code:    code,
		Message: getMessageByCode(code),
		Data:    data,
	}
}

func BizResponse(c *gin.Context, httpCode int, bizCode int32, data interface{}) {
	resp := NewGinResponse(bizCode, data)
	c.JSON(httpCode, resp)
}

func Success(c *gin.Context) {
	resp := NewGinResponse(CodeSuccess, nil)
	c.JSON(http.StatusOK, resp)
}
