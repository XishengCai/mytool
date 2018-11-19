package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Code string
	Price int
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	//Migrate the schema
	db.AutoMigrate(&Product{})

	//Read
	var product Product
	db.First(&product, 1)
	db.First(&product, "code = ?", "L1212")

	db.Model(&product).Update("Price", 20000)

	db.Delete(&product)

}