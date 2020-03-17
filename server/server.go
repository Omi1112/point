package server

import (
	"github.com/SeijiOmi/points-service/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Init is initialize server
func Init() {
	r := router()
	r.Run(":9000")
}

func router() *gin.Engine {
	r := gin.Default()

	// https://godoc.org/github.com/gin-gonic/gin#RouterGroup.Use
	r.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
		},
		// 許可したいHTTPリクエストヘッダの一覧
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"X-Requested-With",
			"Origin",
			"X-Csrftoken",
			"Content-Type",
			"Accept",
		},
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"*",
		},
	}))

	p := r.Group("/points")
	{
		p.GET("/:id", controller.Show)
		p.POST("", controller.Create)
	}

	h := r.Group("/sum")
	{
		h.GET("/:id", controller.Sum)
	}

	return r
}
