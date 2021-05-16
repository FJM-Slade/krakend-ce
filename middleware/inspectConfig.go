package middleware

import (
	config "github.com/devopsfaith/krakend/config"
	//"github.com/devopsfaith/krakend/proxy"
	"github.com/gin-gonic/gin"
)

var timee int = 20

func InspectConfig(cfg config.ExtraConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Do things
	}
}
