package pkg

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 使用gorm模块连接mysql数据库 by soapffz

func GetMysqlConnByGorm(username string, password string, host string, port string, dbname string) (db *gorm.DB) {
	// 数据库初始化及连接语句

	// 通过字符串拼接数据库连接语句
	var builder strings.Builder
	builder.Write([]byte(username))
	builder.Write([]byte(":"))
	builder.Write([]byte(password))
	builder.Write([]byte("@tcp("))
	builder.Write([]byte(host))
	builder.Write([]byte(":"))
	builder.Write([]byte(port))
	builder.Write([]byte(")/"))
	builder.Write([]byte(dbname))
	builder.Write([]byte("?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=60s"))

	// fmt.Printf("字符串拼接为：\n" + builder.String())

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       builder.String(), // DSN data source name"
		DefaultStringSize:         256,              // string 类型字段的默认长度
		DisableDatetimePrecision:  true,             // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,             // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,             // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,            // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，必须设置，不然会在表名后加s
		},
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 1 * time.Millisecond,
			LogLevel:      logger.Error,
			Colorful:      true,
		}),
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)  //空闲连接数
	sqlDB.SetMaxOpenConns(100) //最大连接数,不设置go并发会提示too many connections
	sqlDB.SetConnMaxLifetime(time.Minute)
	if err != nil {
		panic("连接数据库失败！")
	}
	return db
}
