package db

import (
	"final_exam/model"
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func DB() {
	viper.SetConfigName("config")
	viper.AddConfigPath("db")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mysql.dsn"))
	dsn := viper.GetString("mysql.dsn")

	dialactor := mysql.Open(dsn)
	db, err := gorm.Open(dialactor)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection successful")

	Database = db
	// 🔹 ดึงลูกค้าทั้งหมด
	var customers []model.Customer
	Database.Find(&customers)

	// 🔹 วนลูปตรวจสอบและเข้ารหัสรหัสผ่าน (ถ้ายังไม่เข้ารหัส)
	for _, customer := range customers {
		if _, err := bcrypt.Cost([]byte(customer.Password)); err != nil {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
			Database.Model(&customer).Update("password", string(hashedPassword))
		}
	}
}

func GetCustomer(email, password string) ([]model.Customer, error) {
	var customers []model.Customer

	result := Database.Where("email = ? AND password = ?", email, password).Find(&customers)
	if result.Error != nil {
		return nil, result.Error
	}

	return customers, nil
}

func GetCustomerByEmail(email string) (*model.Customer, error) {
	var customer model.Customer
	result := Database.Where("email = ?", email).First(&customer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &customer, nil
}

func UpdatePassword(email, newPassword string) error {
	return Database.Model(&model.Customer{}).Where("email = ?", email).Update("password", newPassword).Error
}

func GetProducts(description string, priceMin, priceMax float64) ([]model.Product, error) {
	var products []model.Product
	query := Database.Table("product")

	// ค้นหาจาก description หากมี
	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}

	// ค้นหาจากราคา หากมี
	if priceMin > 0 {
		query = query.Where("price >= ?", priceMin)
	}
	if priceMax > 0 {
		query = query.Where("price <= ?", priceMax)
	}

	result := query.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

// ฟังก์ชันค้นหารถเข็นของลูกค้า
func GetCartByName(customerID int, cartName string) (model.Cart, error) {
	var cart model.Cart
	result := Database.Where("customer_id = ? AND cart_name = ?", customerID, cartName).First(&cart)
	if result.Error != nil {
		return cart, result.Error
	}
	return cart, nil
}

// ฟังก์ชันเพิ่มสินค้าในรถเข็น
func AddProductToCart(cartID, productID, quantity int) error {
	newCartItem := model.CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
	}
	if err := Database.Create(&newCartItem).Error; err != nil {
		return err
	}
	return nil
}

// ฟังก์ชันอัปเดตจำนวนสินค้าที่มีในรถเข็น
func UpdateProductQuantity(cartItemID, quantity int) error {
	var cartItem model.CartItem
	result := Database.Where("cart_item_id = ?", cartItemID).First(&cartItem)
	if result.Error != nil {
		return result.Error
	}
	cartItem.Quantity += quantity
	if err := Database.Save(&cartItem).Error; err != nil {
		return err
	}
	return nil
}
