package mysql

import (
	"context"
	"fmt"
	"payment/internal/model"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"payment/pkg/config"
	"os"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Order{},
		&model.Payway{},
		&model.User{},
	)
}

func Init(cfg *config.DBConfig) (*gorm.DB, error) {
	
	dsn := fmt.Sprintf(cfg.DSN, 
		  os.Getenv("PAYMENT_USER"), 
	      os.Getenv("PAYMENT_PASSWORD"),
		  os.Getenv("PAYMENT_HOST"),
		  os.Getenv("PAYMENT_PORT"),
		  os.Getenv("PAYMENT_DATABASE"))
	//初始化数据库
	db , err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err!= nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)

	ctx , cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()
	if err := sqlDB.PingContext(ctx);err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("数据库连接失败：%v", err)
	}

	// 自动迁移
	if err := AutoMigrate(db);err!=nil {
		return nil, err
	}
	return db, nil
}

