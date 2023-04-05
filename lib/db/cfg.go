package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cfg struct {
	// 数据库类型
	Type Type `json:"type"`
	// 数据库连接DSN
	Dsn  string   `json:"dsn"`
	Pool *CfgPool `json:"pool,omitempty"`
}

type CfgPool struct {
	MaxIdle     int `json:"maxIdle,omitempty"`     // 空闲连接池中的最大连接数
	MaxOpen     int `json:"maxOpen,omitempty"`     // 数据库的最大打开连接数
	MaxLifeTime int `json:"maxLifeTime,omitempty"` // 连接可以复用的最长时间
}

type Type string

const (
	SQLite Type = "sqlite"

	//MySQL  Type = "mysql"
	//PgSql  Type = "postgresql"
	//TiDB   Type = "tidb"
)

func (customDbType Type) dialector(dsn string) gorm.Dialector {
	switch customDbType {
	case SQLite:
		return sqlite.Open(dsn)
	default:
		return nil
	}
}
