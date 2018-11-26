package reviewer

import (
	"fmt"
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type CreateDatabaseReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode     *ast.CreateDatabaseStmt
	ReviewConfig *config.ReviewConfig
	DBConfig     *config.DBConfig
}

func (this *CreateDatabaseReviewer) Init() {
	this.ReviewMSG = NewReivewMSG()
}

func (this *CreateDatabaseReviewer) Review() *ReviewMSG {
	this.Init()

	if !this.ReviewConfig.RuleAllowCreateDatabase {
		msg := fmt.Sprintf("不允许创建数据库: %v", this.StmtNode.Name)
		this.ReviewMSG.AppendMSG(true, msg)
		return this.ReviewMSG
	}

	// 检测名称长度
	haveError := this.DetectDBNameLength()
	if haveError {
		return this.ReviewMSG
	}

	// 检测命名规则
	haveError = this.DetectDBNameReg()
	if haveError {
		return this.ReviewMSG
	}

	// 检测创建数据库其他选项
	haveError = this.DetectDBOptions()
	if haveError {
		return this.ReviewMSG
	}

	// 检测需要创建的数据库是否在目标实例中已经有
	haveError = this.DetectInstanceDatabase()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

// 检测数据库名长度
func (this *CreateDatabaseReviewer) DetectDBNameLength() (haveError bool) {
	var msg string
	haveError, msg = DetectNameLength(this.StmtNode.Name, this.ReviewConfig.RuleNameLength)
	if haveError {
		msg = fmt.Sprintf("%v. 数据库: %v", msg, this.StmtNode.Name)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}
	return
}

// 检测数据库命名规范
func (this *CreateDatabaseReviewer) DetectDBNameReg() (haveError bool) {
	var msg string
	haveError, msg = DetectNameReg(this.StmtNode.Name, this.ReviewConfig.RuleNameReg)
	if haveError {
		msg = fmt.Sprintf("%v. 数据库: %v", msg, this.StmtNode.Name)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}
	return
}

// 检测创建数据库其他选项值
func (this *CreateDatabaseReviewer) DetectDBOptions() (haveError bool) {
	for _, option := range this.StmtNode.Options {
		var msg string

		switch option.Tp {
		case ast.DatabaseOptionCharset:
			haveError, msg = DetectCharset(option.Value, this.ReviewConfig.RuleCharSet)
		case ast.DatabaseOptionCollate:
			haveError, msg = DetectCollate(option.Value, this.ReviewConfig.RuleCollate)
		}

		// 一检测到有问题键停止下面检测, 返回检测错误值
		if haveError {
			this.ReviewMSG.AppendMSG(haveError, msg)
			break
		}
	}

	return
}

// 链接到实例检测相关信息
func (this *CreateDatabaseReviewer) DetectInstanceDatabase() (haveError bool) {
	var msg string
	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	tableInfo.DBName = this.StmtNode.Name
	err := tableInfo.OpenInstance()
	if err != nil {
		msg = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		this.ReviewMSG.AppendMSG(haveError, msg)
		tableInfo.CloseInstance()
		return
	}

	// 数据库存在报错
	haveError, msg = DetectDatabaseExistsByName(tableInfo, this.StmtNode.Name)
	haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		tableInfo.CloseInstance()
		return
	}

	tableInfo.CloseInstance()
	return
}
