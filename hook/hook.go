package hook

import "github.com/gin-gonic/gin"

type Hooker interface {
	Hook(c Callback) error
}

type Builder interface {
	NewHooker(c *gin.Context) (Hooker, error)
}
