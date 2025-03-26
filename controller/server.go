package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	//Run Release
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is now working",
		})
	})
	//Load Controllers
	CustomerController(router)
	ProductController(router)
	CartController(router)

	router.Run()
}
