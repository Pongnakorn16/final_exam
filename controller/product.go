package controller

import (
	"final_exam/db"
	"final_exam/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductController(router *gin.Engine) {

	routes := router.Group("/product")
	{
		routes.POST("search", get_search)
	}
}

func get_search(c *gin.Context) {
	var req dto.SearchProduct

	// อ่าน JSON และตรวจสอบค่าที่รับมา
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// เรียกใช้ GetProducts เพื่อค้นหาจาก description และ price
	products, err := db.GetProducts(req.Description, req.PriceMin, req.PriceMax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่ง JSON response
	c.JSON(http.StatusOK, gin.H{"products": products})
}
