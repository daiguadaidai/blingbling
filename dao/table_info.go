package dao

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/parser"
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/dependency/mysql"
	"strconv"
)

type TableInfo struct {
	DBName string
	TableName string
	Instance *Instance

	Exists bool
	ExistsQueried bool // 已经查询过了 表是否存在
	DBExists bool
	DBExistsQueried bool // 已经查询过了 数据库是否存在

	CreateTableSql string
	ColumnNameMap map[string]bool
	PKColumnNameMap map[string]bool
	PKColumnNameList []string
	UniqueIndexes map[string][]string
	Indexes map[string][]string
	FullTextIndex map[string][]string
	PartitionNames map[string]bool
	ColumnTypeCount map[byte]int // 保存字段类型出现的个数
}

/* 新建一个表信息
Params:
    _instance: 实例
    _table: 表名
 */
func NewTableInfo(_dbConfig *config.DBConfig, _table string) *TableInfo {
	return &TableInfo{
		DBName: _dbConfig.Database,
		TableName: _table,
		Instance: NewInstance(_dbConfig),
	}
}

// 打开实例链接
func (this *TableInfo) OpenInstance() error {
	return this.Instance.OpenDB()
}

// 关闭实例链接
func (this *TableInfo) CloseInstance() error {
	return this.Instance.CloseDB()
}

/* 检测数据库是否存在
Params:
    _dbName: 数据库名
 */
func (this *TableInfo) DatabaseExistsByName(_dbName string) (bool, error) {
	if this.DBExistsQueried {
		return this.DBExists, nil
	}

	sql := `
    SELECT COUNT(*)
    FROM information_schema.SCHEMATA
    WHERE SCHEMA_NAME = ?;
    `

	var count int
	err := this.Instance.DB.QueryRow(sql, _dbName).Scan(&count)
	if count > 0 {
		this.DBExists = true
	}

	this.DBExistsQueried = true

	return this.DBExists, err
}

/* 通过表名确认表是否存在
Params:
    _dbName: 数据库名称
    _tableName: 表名称
 */
func (this *TableInfo) TableExistsByName(_dbName, _tableName string) (bool, error) {
	sql := `
        SELECT COUNT(*)
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = ?
            AND TABLE_NAME = ?
            AND TABLE_TYPE = 'BASE TABLE'
    `

	var count int
	var exists bool
	err := this.Instance.DB.QueryRow(sql, _dbName, _tableName).Scan(&count)
	if count > 0 {
		exists = true
	}

	return exists, err
}

// 获取表所有字段, 并保存到map中
func (this *TableInfo) FindColumnNameMap() (map[string]bool, error) {
	if this.ColumnNameMap != nil { // 存在了就不再次重复查询数据库
		return this.ColumnNameMap, nil
	}

	sql := `
    SELECT COLUMN_NAME
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = ?
        AND TABLE_NAME = ?;
    `

    rows, err := this.Instance.DB.Query(sql, this.DBName, this.TableName)
    if err != nil {
		errMSG := fmt.Sprintf("获取表所有列出错: %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, errors.New(errMSG)
	}
	defer rows.Close()

	columnNameMap := make(map[string]bool)
	var name string
	for rows.Next() {
		rows.Scan(&name)
		columnNameMap[name] = true
	}

	err = rows.Err()
	if err != nil {
		errMSG := fmt.Sprintf("获取表所有列出错(scan): %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, errors.New(errMSG)
	}

	this.ColumnNameMap = make(map[string]bool)
	this.ColumnNameMap = columnNameMap

	return this.ColumnNameMap, nil
}

// 获取表的所有索引 已经约束
func (this *TableInfo) FindAllIndexes() (map[string][]string, error) {
	// 如果已经有则直接返回, 不必在到数据库中获取
	if this.Indexes != nil {
		return this.Indexes, nil
	}

	// 获取主键
	_, _, err := this.FindPrimaryKey()
	if err != nil {
		return nil, err
	}

	// 获取唯一键
	_, err = this.FindUniqueIndexes()
	if err != nil {
		return nil, err
	}

	// 获取索引
	_, err = this.FindNormalIndexes()
	if err != nil {
		return nil, err
	}

	// 键主键, 都加入到
	this.PrimaryCombinIndexes()
	// 唯一键加入 索引
	this.UniqueIndexCombinIndexes()

	return this.Indexes, nil
}

// 获取主键
func (this *TableInfo) FindPrimaryKey() ([]string, map[string]bool, error) {
	if this.PKColumnNameList != nil && this.PKColumnNameMap != nil {
		return this.PKColumnNameList, this.PKColumnNameMap, nil
	}

	sql := `
        SELECT
            S.COLUMN_NAME
        FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC
        LEFT JOIN INFORMATION_SCHEMA.STATISTICS AS S
            ON TC.TABLE_SCHEMA = S.INDEX_SCHEMA
            AND TC.TABLE_NAME = S.TABLE_NAME
            AND TC.CONSTRAINT_NAME = S.INDEX_NAME 
        WHERE TC.TABLE_SCHEMA = ?
            AND TC.TABLE_NAME = ?
            AND TC.CONSTRAINT_TYPE = 'PRIMARY KEY'
        ORDER BY TC.CONSTRAINT_NAME ASC, S.SEQ_IN_INDEX ASC
    `

	rows, err := this.Instance.DB.Query(sql, this.DBName, this.TableName)
	if err != nil {
		errMSG := fmt.Sprintf("获取表主键列名: %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, nil, errors.New(errMSG)
	}
	defer rows.Close()


	pkColumnNameMap := make(map[string]bool)
	pkColumnNameList := make([]string, 0, 1)
	var pkColumnName string
	for rows.Next() {
		rows.Scan(&pkColumnName)
		pkColumnNameMap[pkColumnName] = true
		pkColumnNameList = append(pkColumnNameList, pkColumnName)
	}

	err = rows.Err()
	if err != nil {
		errMSG := fmt.Sprintf("获取表主键列名(scan): %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, nil, errors.New(errMSG)
	}

	this.PKColumnNameMap = pkColumnNameMap
	this.PKColumnNameList = pkColumnNameList

	return this.PKColumnNameList, this.PKColumnNameMap, nil
}

// 获取唯一约束据信息
func (this *TableInfo) FindUniqueIndexes() (map[string][]string, error) {
	if this.UniqueIndexes != nil {
		return this.UniqueIndexes, nil
	}

	sql := `
        SELECT
            TC.CONSTRAINT_NAME,
            S.COLUMN_NAME
        FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC
        LEFT JOIN INFORMATION_SCHEMA.STATISTICS AS S
            ON TC.TABLE_SCHEMA = S.INDEX_SCHEMA
            AND TC.TABLE_NAME = S.TABLE_NAME
            AND TC.CONSTRAINT_NAME = S.INDEX_NAME 
        WHERE TC.TABLE_SCHEMA = ?
            AND TC.TABLE_NAME = ?
            AND TC.CONSTRAINT_TYPE = 'UNIQUE'
        ORDER BY TC.CONSTRAINT_NAME ASC, S.SEQ_IN_INDEX ASC
    `

	rows, err := this.Instance.DB.Query(sql, this.DBName, this.TableName)
	if err != nil {
		errMSG := fmt.Sprintf("获取唯一键列名: %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, errors.New(errMSG)
	}
	defer rows.Close()


	uniqueIndexes := make(map[string][]string)
	var uniqueName string
	var uniqueColumnName string
	for rows.Next() {
		rows.Scan(&uniqueName, &uniqueColumnName)
		if _, ok := uniqueIndexes[uniqueName]; !ok {
			uniqueIndexes[uniqueName] = make([]string, 0, 1)
		}
		uniqueIndexes[uniqueName] = append(uniqueIndexes[uniqueName], uniqueColumnName)
	}

	err = rows.Err()
	if err != nil {
		errMSG := fmt.Sprintf("获取唯一键列名(scan): %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, errors.New(errMSG)
	}

	this.UniqueIndexes = uniqueIndexes

	return this.UniqueIndexes, nil
}

func (this *TableInfo) FindNormalIndexes() (map[string][]string, error) {
	if this.Indexes != nil {
		return this.Indexes, nil
	}

	sql := fmt.Sprintf("SHOW INDEX FROM `%v`.`%v` WHERE Non_unique = 1",
		this.DBName, this.TableName)

	rows, err := this.Instance.DB.Query(sql)
	if err != nil {
		errMSG := fmt.Sprintf("获取普通索引: %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, errors.New(errMSG)
	}
	defer rows.Close()


	indexes := make(map[string]map[int]string)
	var ignore interface{}
	var indexName string
	var seqInIndex int
	var indexColumnName string
	for rows.Next() {
		rows.Scan(&ignore, &ignore, &indexName, &seqInIndex, &indexColumnName,&ignore,
			&ignore, &ignore, &ignore, &ignore, &ignore, &ignore, &ignore)
		if _, ok := indexes[indexName]; !ok {
			indexes[indexName] = make(map[int]string)
		}

		indexes[indexName][seqInIndex] = indexColumnName
	}

	err = rows.Err()
	if err != nil {
		errMSG := fmt.Sprintf("获取普通索引(scan): %v.%v %v:%v. %v",
			this.DBName, this.TableName, this.Instance.DBconfig.Host, this.Instance.DBconfig.Port,
			err)
		return nil, errors.New(errMSG)
	}

	// 将查询主来的转化成想要的格式
	this.Indexes = make(map[string][]string)
	for indexName, indexColumnNameMap := range indexes {
		if _, ok := this.Indexes[indexName]; !ok {
			this.Indexes[indexName] = make([]string, len(indexColumnNameMap))
		}

		for seqInIndex, columnName := range indexColumnNameMap {
			this.Indexes[indexName][seqInIndex - 1] = columnName
		}
	}

	return this.Indexes, nil
}

// 将主键加入索引中
func (this *TableInfo) PrimaryCombinIndexes() {
	if this.Indexes == nil {
		this.Indexes = make(map[string][]string)
	}

	// 键主键加入索引中
	if this.PKColumnNameList != nil {
		this.Indexes["PRIMARY"] = this.PKColumnNameList
	}
}

// 将唯一索引加入索引中
func (this *TableInfo) UniqueIndexCombinIndexes() {
	if this.Indexes == nil {
		this.Indexes = make(map[string][]string)
	}

	// 将唯一索引加入索引中
	if this.UniqueIndexes != nil {
		for uniqueName, uniqueIndex := range this.UniqueIndexes {
			this.Indexes[uniqueName] = uniqueIndex
		}
	}
}

/* 获取表的建表sql
Params:
    _dbName: 数据库名称
    _tableName: 表名
 */
func (this *TableInfo) InitCreateTableSql(_dbName, _tableName string) error {
	if _dbName == "" {
		_dbName = this.Instance.DBconfig.Database
	}

	sql := fmt.Sprintf("SHOW CREATE TABLE `%v`.`%v`", _dbName, _tableName)

	var ignore string
	err := this.Instance.DB.QueryRow(sql).Scan(&ignore, &this.CreateTableSql)
	if err != nil {
		errMSG := fmt.Sprintf("show create table `%v`.`%v` 失败. %v",
			_dbName, _tableName, err)
		return  errors.New(errMSG)
	}

	return nil
}

func (this *TableInfo) GetExplainMaxRows(_sql string) (int, error) {
	var maxRowCount int
	var rowsColumnIndex int

	rows, err := this.Instance.DB.Query(_sql)
	if err != nil {
		errMSG := fmt.Sprintf("执行explain失败: %v:%v. %v %v",
		this.Instance.DBconfig.Host, this.Instance.DBconfig.Port, _sql, err)
		return -1, errors.New(errMSG)
	}
	defer rows.Close()

	// 获取rows 字段在第几个
	columnNames, err := rows.Columns()
	if err != nil {
		errMSG := fmt.Sprintf("获取字段错误: %v:%v. %v %v",
			this.Instance.DBconfig.Host, this.Instance.DBconfig.Port, _sql, err)
		return -1, errors.New(errMSG)
	}
	for i, columnName := range columnNames {
		if columnName == "rows" {
			rowsColumnIndex = i
		}
	}

	scanArgs := make([]interface{}, len(columnNames))
	values := make([]interface{}, len(columnNames))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		var rowCount int
		if values[rowsColumnIndex] != nil {
			rowCount, err = strconv.Atoi(string(values[rowsColumnIndex].([]uint8)))
			if err != nil {
				errMSG := fmt.Sprintf("explain rows值转化为数字错误: %v:%v. %v %v",
					this.Instance.DBconfig.Host, this.Instance.DBconfig.Port, _sql, err)
				return -1, errors.New(errMSG)
			}
		}
		if maxRowCount < rowCount {
			maxRowCount = rowCount
		}
	}

	err = rows.Err()
	if err != nil {
		errMSG := fmt.Sprintf("执行explain失败(scan): %v:%v. %v %v",
			this.Instance.DBconfig.Host, this.Instance.DBconfig.Port, _sql, err)
		return -1, errors.New(errMSG)
	}

	return maxRowCount, nil
}

// 解析表建表语句获取相关信息
func (this *TableInfo) ParseCreateTableInfo() error {
	// 解析SQL
	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(this.CreateTableSql, "", "")
	if err != nil {
		errMSG := fmt.Sprintf("解析原实例建表sql语法错误: %v", err)
		return errors.New(errMSG)
	}

	for _, stmtNode := range stmtNodes {
		createTableStmt := stmtNode.(*ast.CreateTableStmt)

		this.ParseCreateTableColumns(createTableStmt)
		this.ParseCreateTableConstraint(createTableStmt)
		this.ParseCreateTablePartition(createTableStmt)

		break
	}
	return nil
}

/* 解析创建表
Params:
    _createTableStmt: 创建表语句
 */
func (this *TableInfo) ParseCreateTableColumns(_createTableStmt *ast.CreateTableStmt) {
	if this.ColumnNameMap == nil {
		this.ColumnNameMap = make(map[string]bool)
	}

	for _, column := range _createTableStmt.Cols {
		this.ColumnNameMap[column.Name.String()] = true
	}
}

/* 解析建表约束
Params:
    _createTableStmt: 创建表语句
 */
func (this *TableInfo) ParseCreateTableConstraint(_createTableStmt *ast.CreateTableStmt) {
	if this.PKColumnNameList == nil {
		this.PKColumnNameList = make([]string, 0, 1)
	}
	if this.PKColumnNameMap == nil {
		this.PKColumnNameMap = make(map[string]bool)
	}
	if this.Indexes == nil {
		this.Indexes = make(map[string][]string)
	}
	if this.UniqueIndexes == nil {
		this.UniqueIndexes = make(map[string][]string)
	}
	if this.FullTextIndex == nil {
		this.FullTextIndex = make(map[string][]string)
	}

	for _, constraint := range _createTableStmt.Constraints {
		// 获取索引列名列表
		indexColumns := make([]string, 0, 1)
		for _, columnName := range constraint.Keys {
			indexColumns = append(indexColumns, columnName.Column.String())
		}

		switch constraint.Tp {
		case ast.ConstraintNoConstraint:
		case ast.ConstraintPrimaryKey:
			for _, columnName := range indexColumns {
				this.PKColumnNameMap[columnName] = true
			}
			this.PKColumnNameList = indexColumns
			this.UniqueIndexes["PRIMARY KEY"] = indexColumns
			this.Indexes["PRIMARY KEY"] = indexColumns
		case ast.ConstraintKey, ast.ConstraintIndex:
			this.Indexes[constraint.Name] = indexColumns
		case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
			this.UniqueIndexes[constraint.Name] = indexColumns
			this.Indexes[constraint.Name] = indexColumns
		case ast.ConstraintForeignKey:
		case ast.ConstraintFulltext:
			this.Indexes[constraint.Name] = indexColumns
			this.Indexes[constraint.Name] = indexColumns
		}
	}
}

/* 解析建表语句的分区表
    _createTableStmt: 创建表语句
 */
func (this *TableInfo) ParseCreateTablePartition(_createTableStmt *ast.CreateTableStmt) {
	if this.PartitionNames == nil {
		this.PartitionNames = make(map[string]bool)
	}

	if _createTableStmt.Partition != nil {
		for _, partition := range _createTableStmt.Partition.Definitions {
			this.PartitionNames[partition.Name.String()] = true
		}
	}
}

/* 解析键表语句的表字段个数
Params:
    _createTableStmt: 建表语句
 */
func (this *TableInfo) ParseCreateTableColumnTableCount(_createTableStmt *ast.CreateTableStmt) {
	for _, column := range _createTableStmt.Cols {
		switch column.Tp.Tp {
		case mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeBlob:
			// 4种大字段都设置为是 Blob
			this.ColumnTypeCount[mysql.TypeBlob] ++
		default:
			this.ColumnTypeCount[column.Tp.Tp] ++
		}
	}
}
