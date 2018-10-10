package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type CreateDatabaseReviewer struct {
	StmtNode *ast.CreateDatabaseStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
}

func (this *CreateDatabaseReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowCreateDatabase {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("不允许创建数据库: %v", this.StmtNode.Name)
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测名称长度
	reviewMSG = this.DetectDBNameLength()
	if reviewMSG != nil {
		reviewMSG.MSG = fmt.Sprintf("%v %v", "数据库名", reviewMSG.MSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测命名规则
	reviewMSG = this.DetectDBNameReg()
	if reviewMSG != nil {
		reviewMSG.MSG = fmt.Sprintf("%v %v", "数据库名", reviewMSG.MSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测创建数据库其他选项
	reviewMSG = this.DetectDBOptions()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测需要创建的数据库是否在目标实例中已经有
	reviewMSG = this.DetectInstanceDatabase()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 能走到这里说明写的语句审核成功
	reviewMSG = new(ReviewMSG)
	reviewMSG.Code = REVIEW_CODE_SUCCESS
	reviewMSG.MSG = "审核成功"

	return reviewMSG
}

// 检测数据库名长度
func (this *CreateDatabaseReviewer) DetectDBNameLength() *ReviewMSG {
	return DetectNameLength(this.StmtNode.Name, this.ReviewConfig.RuleNameLength)
}

// 检测数据库命名规范
func (this *CreateDatabaseReviewer) DetectDBNameReg() *ReviewMSG {
	return DetectNameReg(this.StmtNode.Name, this.ReviewConfig.RuleNameReg)
}

// 检测创建数据库其他选项值
func (this *CreateDatabaseReviewer) DetectDBOptions() *ReviewMSG {
	var reviewMSG *ReviewMSG

	for _, option := range this.StmtNode.Options {
		switch option.Tp {
		case ast.DatabaseOptionCharset:
			reviewMSG = DetectCharset(option.Value, this.ReviewConfig.RuleCharSet)
		case ast.DatabaseOptionCollate:
			reviewMSG = DetectCollate(option.Value, this.ReviewConfig.RuleCollate)
		}

		// 一检测到有问题键停止下面检测, 返回检测错误值
		if reviewMSG != nil {
			break
		}
	}

	return reviewMSG
}

// 链接到实例检测相关信息
func (this *CreateDatabaseReviewer) DetectInstanceDatabase() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	tableInfo.DBName = this.StmtNode.Name
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		return reviewMSG
	}

	// 数据库存在报错
	reviewMSG = DetectDatabaseExistsByName(tableInfo, this.StmtNode.Name)
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}
