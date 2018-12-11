package reviewer

import (
	"fmt"

	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
)

type DropDatabaseReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode     *ast.DropDatabaseStmt
	ReviewConfig *config.ReviewConfig
	DBConfig     *config.DBConfig
}

func (this *DropDatabaseReviewer) Init() {
	this.ReviewMSG = NewReivewMSG()
}

func (this *DropDatabaseReviewer) Review() *ReviewMSG {
	this.Init()

	if !this.ReviewConfig.RuleAllowDropDatabase {
		msg := config.MSG_ALLOW_DROP_DATABASE_ERROR
		this.ReviewMSG.AppendMSG(true, msg)
		return this.ReviewMSG
	}

	// 链接数据库检测实例相关信息
	haveError := this.DetectInstanceDatabase()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

// 链接到实例检测相关信息
func (this *DropDatabaseReviewer) DetectInstanceDatabase() (haveError bool) {
	var msg string

	tableInfo := NewTableInfo(this.DBConfig, "")
	tableInfo.DBName = this.StmtNode.Name
	err := tableInfo.OpenInstance()
	if err != nil {
		msg = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		tableInfo.CloseInstance()
		this.ReviewMSG.AppendMSG(false, msg)
		return
	}

	// 数据库不错在报错
	haveError, msg = DetectDatabaseNotExistsByName(tableInfo, this.StmtNode.Name)
	haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		tableInfo.CloseInstance()
		return
	}

	tableInfo.CloseInstance()
	return
}
