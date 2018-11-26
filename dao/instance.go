package dao

import (
	"database/sql"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/daiguadaidai/blingbling/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
)

type Instance struct {
	DBconfig *config.DBConfig
	DB       *sql.DB
}

/* 新建一个数据库执行器
Params:
    _dbConfig: 数据库配置
 */
func NewInstance(_dbConfig *config.DBConfig) *Instance {
	executor := &Instance{DBconfig: _dbConfig}

	return executor
}

// 打开数据库连接
func (this *Instance) OpenDB() error {
	var err error

	this.DB, err = sql.Open("mysql", this.DBconfig.GetDataSource())
	if err != nil { // 打开数据库失败
		errMSG := fmt.Sprintf("打开数据库链接失败[%v:%v]",
			this.DBconfig.Host, this.DBconfig.Port)
		return errors.New(errMSG)
	}

	return err
}

// 关闭数据库链接
func (this *Instance) CloseDB() error {
	var err error

	if this.DB != nil {
		err = this.DB.Close()
		if err != nil {
			log.Errorf("警告: 链接实例检测表相关信息. 关闭连接出错 %v:%v/%v",
				this.DBconfig.Host,
				this.DBconfig.Port,
				this.DBconfig.Database)
		}
	}

	return err
}
