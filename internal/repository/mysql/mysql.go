package mysql

import (
	"context"
	"fmt"
	"payment/internal/model"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"net"
)

type MySqlDB struct {
	DSN string `yaml:"dsn"`
	Host string `yaml:"host"`
	DBName string `yaml:"dbname"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Port string `yaml:"port"`
}

type PaymentDB struct {
	*MySqlDB `yaml:"payment"` 	
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Order{},
		&model.Payway{},
		&model.User{},
	)
}

func Init(cfg *MySqlDB) (*gorm.DB, error) {
	// 注册自定义 Dialer，添加底层 TCP 连接超时
	mysqlDriver.RegisterDialContext("timeoutDialer", func(ctx context.Context, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout: 3 * time.Second,
		}
		return dialer.DialContext(ctx, "tcp", addr)
	})
	
	dsn := fmt.Sprintf(cfg.DSN, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
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

