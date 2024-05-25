package initdb

import (
	"fmt"
	"ktserver/pkg/db"
	"os"
)

type InitDBModel struct {
	DBType   string `json:"dbType"`                    // 数据库类型
	Host     string `json:"host"`                      // 服务器地址
	Port     string `json:"port"`                      // 数据库连接端口
	UserName string `json:"userName"`                  // 数据库用户名
	Password string `json:"password"`                  // 数据库密码
	DBName   string `json:"dbName" binding:"required"` // 数据库名
	DBPath   string `json:"dbPath"`                    // sqlite数据库文件路径
}

// MysqlEmptyDsn msyql 空数据库 建库链接
// Author SliverHorn
func (i *InitDBModel) MysqlEmptyDsn() string {
	if i.Host == "" {
		i.Host = "127.0.0.1"
	}
	if i.Port == "" {
		i.Port = "3306"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", i.UserName, i.Password, i.Host, i.Port)
}

// PgsqlEmptyDsn pgsql 空数据库 建库链接
// Author SliverHorn
func (i *InitDBModel) PgsqlEmptyDsn() string {
	if i.Host == "" {
		i.Host = "127.0.0.1"
	}
	if i.Port == "" {
		i.Port = "5432"
	}
	return "host=" + i.Host + " user=" + i.UserName + " password=" + i.Password + " port=" + i.Port + " dbname=" + "postgres" + " " + "sslmode=disable TimeZone=Asia/Shanghai"
}

// SqliteEmptyDsn sqlite 空数据库 建库链接
// Author Kafumio
func (i *InitDBModel) SqliteEmptyDsn() string {
	separator := string(os.PathSeparator)
	return i.DBPath + separator + i.DBName + ".db"
}

func (i *InitDBModel) MssqlEmptyDsn() string {
	return "sqlserver://" + i.UserName + ":" + i.Password + "@" + i.Host + ":" + i.Port + "?database=" + i.DBName + "&encrypt=disable"
}

// ToMysqlConfig 转换 config.Mysql
// Author [SliverHorn](https://github.com/SliverHorn)
func (i *InitDBModel) ToMysqlConfig() *db.MySQLOptions {
	return &db.MySQLOptions{
		Username: i.UserName,
		Password: i.Password,
		Host:     i.Host + ":" + i.Port,
		Database: i.DBName,
	}
}
