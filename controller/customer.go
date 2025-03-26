package controller

import (
	"final_exam/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CustomerController(router *gin.Engine) {

	router.POST("/auth/login", get_login)
}

func get_login(c *gin.Context) {
	var req LoginRequest

	// อ่าน JSON และตรวจสอบค่าที่รับมา
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// เรียกใช้ GetCustomer ด้วย email และ password ที่รับจาก JSON
	customers, err := db.GetCustomer(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "ไม่มีข้อมูลลูกค้า"})
		return
	}

	// ส่ง JSON response
	c.JSON(http.StatusOK, gin.H{"customers": customers})
}
