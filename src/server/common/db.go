package common


import (
	"github.com/hoisie/mustache"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	_ "github.com/go-sql-driver/mysql"
)

func _SetupDB(name string, readOnly bool) *gorm.DB {
	isDevelopment := viper.GetBool("isDevelopment")
	instance := "wr"
	if readOnly {
		instance = "r"
	}
	prefix := mustache.Render("resource.database.{{name}}.{{instance}}", map[string]interface{}{
		"name":     name,
		"instance": instance,
	})
	dialect := viper.GetString(prefix + ".dialect")
	user := viper.GetString(prefix + ".user")
	password := viper.GetString(prefix + ".password")
	database := viper.GetString(prefix + ".database")
	host := viper.GetString(prefix + ".host")
	port := viper.GetString(prefix + ".port")
	maxIdle := viper.GetInt(prefix + ".max-idle")
	maxOpen := viper.GetInt(prefix + ".max-open")
	url := mustache.Render("{{user}}:{{password}}@tcp({{host}}:{{port}})/{{database}}?charset=utf8&parseTime=True&loc=Local", map[string]interface{}{
		"user":     user,
		"password": password,
		"database": database,
		"host":     host,
		"port":     port,
	})
	db, err := gorm.Open(dialect, url)
	if err != nil {
		Logger.Warningln("failed to connect database:", database, err.Error())
		panic("failed to connect database")
	}
	db.LogMode(isDevelopment)
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	return db
}

func SetupDB() {

	SodaDB_WR = _SetupDB("soda", false)
	SodaDB_R = _SetupDB("soda", true)

	SodaMngDB_WR = _SetupDB("soda-manager", false)
	SodaMngDB_R = _SetupDB("soda-manager", true)

}

var (
	SodaDB_WR    *gorm.DB
	SodaDB_R     *gorm.DB
	SodaMngDB_WR *gorm.DB
	SodaMngDB_R  *gorm.DB
)
