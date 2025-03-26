// controller/cart_controller.go
package controller

import (
	"final_exam/db"
	"final_exam/dto"
	"final_exam/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartController(router *gin.Engine) {

	routes := router.Group("/cart")
	{
		routes.POST("add", addToCart)
		routes.GET("detail", GetCartDetails)
	}
}

func addToCart(c *gin.Context) {
	var req dto.CartItemRequest

	// อ่านข้อมูล JSON ที่รับมา
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 1. ค้นหารถเข็นที่มีชื่อที่กำหนด
	cart, err := db.GetCartByName(req.CustomerID, req.CartName)
	if err != nil {
		if err.Error() == "record not found" {
			// ถ้าไม่พบรถเข็น ให้สร้างรถเข็นใหม่
			newCart := model.Cart{
				CustomerID: req.CustomerID,
				CartName:   req.CartName,
			}
			if err := db.Database.Create(&newCart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart: " + err.Error()})
				return
			}
			cart = newCart
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 2. ตรวจสอบว่ามีสินค้าในรถเข็นแล้วหรือไม่
	var cartItem model.CartItem
	result := db.Database.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem)

	if result.Error != nil {
		// ถ้าไม่มีสินค้าในรถเข็น ให้เพิ่มสินค้าใหม่
		if err := db.AddProductToCart(cart.CartID, req.ProductID, req.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
	} else {
		// ถ้ามีสินค้าแล้ว ให้เพิ่มจำนวน
		if err := db.UpdateProductQuantity(cartItem.CartItemID, req.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product quantity updated"})
	}
}

func GetCartDetails(c *gin.Context) {
	customerID := c.DefaultQuery("customer_id", "0") // ใช้ query parameter สำหรับ customer_id

	// ค้นหารถเข็นทั้งหมดของลูกค้า
	cartDetails, err := db.GetCartDetailsByCustomerID(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart details: " + err.Error()})
		return
	}

	// ส่งข้อมูล cartDetails กลับในรูปแบบ JSON โดยไม่ต้องใช้ DTO
	c.JSON(http.StatusOK, gin.H{"cart_details": cartDetails})
}
