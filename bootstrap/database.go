package bootstrap

import (
	"fmt"
	"go-fiber/core/logs"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// **************************************** Mysql ****************************************************
func NewDatabaseConnection(env *Env) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.Database.Username,
		env.Database.Password,
		env.Database.DBHost,
		env.Database.DBPort,
		env.Database.DBName,
	)
	// logs.Info(dsn)
	// Set up a custom pool configuration
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("can not connect to database")
	}
	logs.Info("database connection success")

	// Migrate(db)

	return db
}

//**************************************** Postgres ****************************************************

// func NewDatabaseConnectionPostgres(env *Env) *gorm.DB {
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		env.Database.DBHost,
// 		env.Database.DBPort,
// 		env.Database.Username,
// 		env.Database.Password,
// 		env.Database.DBName,
// 	)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		panic("can not connect to database")
// 	}
// 	logs.Info("database connection success")

// 	// Migrate(db)

// 	return db
// }

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
	//Migrate entities
	)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		logs.Error(err.Error())
	}
	logs.Info("Migrate successfully")
}
