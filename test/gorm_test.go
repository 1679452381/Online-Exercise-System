package test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"online_exercise_system/models"
	"testing"
)

func TestGorm(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/online_exercise?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err.Error())
	}
	data := make([]*models.User, 0)
	err = db.Find(&data).Error
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, datum := range data {
		fmt.Println(datum)
	}
}
