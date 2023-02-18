package response

import (
	"context"
	"github.com/gin-gonic/gin"
)

func GenCtx(c *gin.Context) context.Context {
	ctx := c.Copy()
	return ctx
}
