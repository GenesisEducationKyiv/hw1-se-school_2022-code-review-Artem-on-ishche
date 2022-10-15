package presentation

import (
	"encoding/json"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"customers-service/config"
	"customers-service/model"
)

func extractCode(err error) int {
	if err == nil {
		return http.StatusOK
	} else {
		return http.StatusInternalServerError
	}
}

func getDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.MySqlDsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&model.Customer{})

	return db
}

func structToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}
