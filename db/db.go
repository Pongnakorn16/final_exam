package db

import (
	"final_exam/model"
	"fmt"

	"github.com/spf13/viper"
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

	// customer := []model.Customer{}

	// result := db.Find(&customer)
	// if result.Error != nil {
	// 	panic(result.Error)
	// }
	// fmt.Printf("%v", customer)
}

func GetCustomer(email, password string) ([]model.Customer, error) {
	var customers []model.Customer

	result := Database.Where("email = ? AND password = ?", email, password).Find(&customers)
	if result.Error != nil {
		return nil, result.Error
	}

	return customers, nil
}
