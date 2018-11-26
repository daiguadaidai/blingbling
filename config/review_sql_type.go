package config

import (
	log "github.com/cihub/seelog"
	"github.com/daiguadaidai/blingbling/dependency/mysql"
	"sync"
)

const (
	TYPE_STR_DECIMAL     = "decimal"
	TYPE_STR_TINYINT     = "tinyint"
	TYPE_STR_SHORT       = "smallint"
	TYPE_STR_LONG        = "int"
	TYPE_STR_FLOAT       = "float"
	TYPE_STR_DOUBLE      = "double"
	TYPE_STR_NULL        = "NULL"
	TYPE_STR_TIMESTAMP   = "timestamp"
	TYPE_STR_LONG_LONG   = "bigint"
	TYPE_STR_INT24       = "mediumint"
	TPPE_STR_DATE        = "date"
	TYPE_STR_DURATION    = "time"
	TYPE_STR_DATETIME    = "datetime"
	TYPE_STR_YEAR        = "year"
	TYPE_STR_NEW_DATE    = "newdate"
	TYPE_STR_VARCHAR     = "varchar"
	TYPE_STR_BIT         = "bit"
	TYPE_STR_JSON        = "json"
	TYPE_STR_NEW_DECIMAL = "newdecimal"
	TYPE_STR_ENUM        = "enum"
	TYPE_STR_SET         = "set"
	TYPE_STR_TINYBLOB    = "tinyblob"
	TYPE_STR_MEDIUMBLOB  = "mediumblob"
	TYPE_STR_LONG_BLOB   = "longblob"
	TYPE_STR_BLOB        = "blob"
	TYPE_STR_TINYTEXT    = "tinytext"
	TYPE_STR_MEDIUMTEXT  = "mediumtext"
	TYPE_STR_LONG_TEXT   = "longtext"
	TYPE_STR_TEXT        = "text"
	TYPE_STR_VARSTRING   = "varstring"
	TYPE_STR_STRING      = "string"
	TYPE_STR_GEOMETRY    = "geometry"
)

var reviewSQLTypeStringByte map[string]byte
var reviewSQLTypeStringByteOnce sync.Once

var reviewSQLTypeByteString map[byte]string
var reviewSQLTypeByteStringOnce sync.Once

// 获取 字段类型(字符) 到 (Byte) 的映射
func GetReviewSQLTypeStringByte() map[string]byte {
	reviewSQLTypeStringByteOnce.Do(func() {
		reviewSQLTypeStringByte = map[string]byte{
			TYPE_STR_DECIMAL:     mysql.TypeDecimal,
			TYPE_STR_TINYINT:     mysql.TypeTiny,
			TYPE_STR_SHORT:       mysql.TypeShort,
			TYPE_STR_LONG:        mysql.TypeLong,
			TYPE_STR_FLOAT:       mysql.TypeFloat,
			TYPE_STR_DOUBLE:      mysql.TypeDouble,
			TYPE_STR_NULL:        mysql.TypeNull,
			TYPE_STR_TIMESTAMP:   mysql.TypeTimestamp,
			TYPE_STR_LONG_LONG:   mysql.TypeLonglong,
			TYPE_STR_INT24:       mysql.TypeInt24,
			TPPE_STR_DATE:        mysql.TypeDate,
			TYPE_STR_DURATION:    mysql.TypeDuration,
			TYPE_STR_DATETIME:    mysql.TypeDatetime,
			TYPE_STR_YEAR:        mysql.TypeYear,
			TYPE_STR_NEW_DATE:    mysql.TypeNewDate,
			TYPE_STR_VARCHAR:     mysql.TypeVarchar,
			TYPE_STR_BIT:         mysql.TypeBit,
			TYPE_STR_JSON:        mysql.TypeJSON,
			TYPE_STR_NEW_DECIMAL: mysql.TypeNewDecimal,
			TYPE_STR_ENUM:        mysql.TypeEnum,
			TYPE_STR_SET:         mysql.TypeSet,
			TYPE_STR_TINYBLOB:    mysql.TypeTinyBlob,
			TYPE_STR_MEDIUMBLOB:  mysql.TypeMediumBlob,
			TYPE_STR_LONG_BLOB:   mysql.TypeLongBlob,
			TYPE_STR_BLOB:        mysql.TypeBlob,
			TYPE_STR_TINYTEXT:    mysql.TypeTinyBlob,
			TYPE_STR_MEDIUMTEXT:  mysql.TypeMediumBlob,
			TYPE_STR_LONG_TEXT:   mysql.TypeLongBlob,
			TYPE_STR_TEXT:        mysql.TypeBlob,
			TYPE_STR_VARSTRING:   mysql.TypeVarchar,
			TYPE_STR_STRING:      mysql.TypeString,
			TYPE_STR_GEOMETRY:    mysql.TypeGeometry,
		}

	})

	return reviewSQLTypeStringByte
}

// 获取 字段类型(Byte) 到 (字符) 的映射
func GetReviewSQLTypeByteString() map[byte]string {
	reviewSQLTypeByteStringOnce.Do(func() {
		reviewSQLTypeByteString = map[byte]string{
			mysql.TypeDecimal:    TYPE_STR_DECIMAL,
			mysql.TypeTiny:       TYPE_STR_TINYINT,
			mysql.TypeShort:      TYPE_STR_SHORT,
			mysql.TypeLong:       TYPE_STR_LONG,
			mysql.TypeFloat:      TYPE_STR_FLOAT,
			mysql.TypeDouble:     TYPE_STR_DOUBLE,
			mysql.TypeNull:       TYPE_STR_NULL,
			mysql.TypeTimestamp:  TYPE_STR_TIMESTAMP,
			mysql.TypeLonglong:   TYPE_STR_LONG_LONG,
			mysql.TypeInt24:      TYPE_STR_INT24,
			mysql.TypeDate:       TPPE_STR_DATE,
			mysql.TypeDuration:   TYPE_STR_DURATION,
			mysql.TypeDatetime:   TYPE_STR_DATETIME,
			mysql.TypeYear:       TYPE_STR_YEAR,
			mysql.TypeNewDate:    TYPE_STR_NEW_DATE,
			mysql.TypeVarchar:    TYPE_STR_VARCHAR,
			mysql.TypeBit:        TYPE_STR_BIT,
			mysql.TypeJSON:       TYPE_STR_JSON,
			mysql.TypeNewDecimal: TYPE_STR_NEW_DECIMAL,
			mysql.TypeEnum:       TYPE_STR_ENUM,
			mysql.TypeSet:        TYPE_STR_SET,
			mysql.TypeTinyBlob:   TYPE_STR_TINYBLOB,
			mysql.TypeMediumBlob: TYPE_STR_MEDIUMBLOB,
			mysql.TypeLongBlob:   TYPE_STR_LONG_BLOB,
			mysql.TypeBlob:       TYPE_STR_BLOB,
			mysql.TypeString:     TYPE_STR_STRING,
			mysql.TypeGeometry:   TYPE_STR_GEOMETRY,
		}

	})

	return reviewSQLTypeByteString
}

// 通过sql type获取对应类型字符
func GetTextSqlTypeByByte(_sqlType byte) string {
	sqlTypeText := ""

	sqlTypeMap := GetReviewSQLTypeByteString()
	sqlTypeText, ok := sqlTypeMap[_sqlType]
	if !ok {
		log.Errorf("通过MySQL类型(byte)获取对应的字符类型(string)失败")
	}

	return sqlTypeText
}
