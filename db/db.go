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

	// ค้นหาสินค้าจากคำอธิบายและช่วงราคา
	result := Database.Where("description LIKE ? AND price BETWEEN ? AND ?", "%"+description+"%", priceMin, priceMax).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}
