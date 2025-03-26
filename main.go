package main

import (
	"final_exam/controller"
	"final_exam/db"
	"fmt"
)

func main() {
	fmt.Println("DOKKUNG")
	db.DB()
	controller.StartServer()
}
