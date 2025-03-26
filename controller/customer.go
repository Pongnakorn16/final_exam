package controller

import (
	"final_exam/db"
	"final_exam/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CustomerController(router *gin.Engine) {

	routes := router.Group("/customer")
	{
		routes.POST("login", get_login)
		routes.POST("edit_pass", edit_pass)
	}
}

func get_login(c *gin.Context) {
	var req dto.LoginRequest

	// อ่าน JSON และตรวจสอบค่าที่รับมา
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// ค้นหาลูกค้าตาม email
	customer, err := db.GetCustomerByEmail(req.Email)
	if err != nil || customer.CustomerID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "ไม่มีข้อมูลลูกค้า"})
		return
	}

	// ตรวจสอบรหัสผ่าน (เปรียบเทียบกับ bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "รหัสผ่านไม่ถูกต้อง"})
		return
	}

	// ส่ง JSON response
	c.JSON(http.StatusOK, gin.H{"message": "เข้าสู่ระบบสำเร็จ", "customer": customer})
}

func edit_pass(c *gin.Context) {
	var req dto.EditPassRequest

	// ตรวจสอบ JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// ค้นหาลูกค้าตาม Email
	customer, err := db.GetCustomerByEmail(req.Email)
	if err != nil || customer.CustomerID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "ไม่มีข้อมูลลูกค้า"})
		return
	}

	// ตรวจสอบรหัสผ่านเก่าด้วย bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "รหัสผ่านเก่าไม่ถูกต้อง"})
		return
	}

	// เข้ารหัสรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการเข้ารหัส"})
		return
	}

	// อัปเดตรหัสผ่านใหม่
	err = db.UpdatePassword(req.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการอัปเดตรหัสผ่าน"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "เปลี่ยนรหัสผ่านสำเร็จ"})
}
