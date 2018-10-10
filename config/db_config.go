package config

import (
	"fmt"
)

type DBConfig struct {
	Username          string
	Password          string
	Database          string
	CharSet           string
	Host              string
	TimeOut           int
	Port              int
	MaxOpenConns      int
	MaxIdelConns      int
	AllowOldPasswords int
	AutoCommit        bool
}

/* 新建一个数据库执行器
Params:
    _host: ip
    _port: 端口
    _username: 链接数据库用户名
    _password: 链接数据库密码
    _database: 要操作的数据库
 */
func NewDBConfig(
	_host string,
	_port int,
	_username string,
	_password string,
	_database string,
) *DBConfig {
	dbConfig := &DBConfig {
		Username:          _username,
		Password:          _password,
		Host:              _host,
		Port:              _port,
		Database:          _database,
		CharSet:           "utf8mb4",
		MaxOpenConns:      1,
		MaxIdelConns:      1,
		TimeOut:           300,
		AllowOldPasswords: 1,
		AutoCommit:        true,
	}

	return dbConfig
}


func (this *DBConfig) GetDataSource() string {
	dataSource := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=%v&allowOldPasswords=%v&timeout=%vs&autocommit=%v&parseTime=True&loc=Local",
		this.Username,
		this.Password,
		this.Host,
		this.Port,
		this.Database,
		this.CharSet,
		this.AllowOldPasswords,
		this.TimeOut,
		this.AutoCommit,
	)

	return dataSource
}
