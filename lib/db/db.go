package db

import (
	"database/sql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	apiDB  *DB
	onceDB sync.Once
)

func Api() *DB {
	onceDB.Do(func() {
		apiDB = new(DB)
	})
	return apiDB
}

type DB struct {
	cfg                   *Cfg
	ptrGormDB             *gorm.DB
	ptrSqlDB              *sql.DB
	cfgSeqValidatorChains []func() error
}

func (customDB *DB) Init(cfg *Cfg) error {
	if customDB.ptrGormDB == nil && customDB.ptrSqlDB == nil {
		return nil
	}

	customDB.cfg = cfg
	if err := customDB.validateCfg(); err != nil {
		return err
	}

	if err := customDB.initGormDB(); err != nil {
		return err
	}

	customDB.initPool()

	return nil
}

func (customDB *DB) GetGormDB() *gorm.DB {
	return customDB.ptrGormDB
}

func (customDB *DB) GetSqlDB() *sql.DB {
	return customDB.ptrSqlDB
}

func (customDB *DB) initGormDB() error {
	if customDB.ptrGormDB != nil {
		return nil
	}

	tmpDb, err := gorm.Open(customDB.cfg.Type.dialector(customDB.cfg.Dsn))
	if err != nil {
		return err
	}
	tmpSqlDb, err := tmpDb.DB()
	if err != nil {
		return err
	}

	customDB.ptrGormDB = tmpDb
	customDB.ptrSqlDB = tmpSqlDb
	return nil
}

func (customDB *DB) initPool() {
	if customDB.ptrGormDB == nil || customDB.ptrSqlDB == nil {
		return
	}
	if customDB.cfg.Pool != nil {
		if customDB.cfg.Pool.MaxIdle > 0 {
			customDB.ptrSqlDB.SetMaxIdleConns(customDB.cfg.Pool.MaxIdle)
		}
		if customDB.cfg.Pool.MaxOpen > 0 {
			customDB.ptrSqlDB.SetMaxOpenConns(customDB.cfg.Pool.MaxOpen)
		}
		if customDB.cfg.Pool.MaxLifeTime > 0 {
			customDB.ptrSqlDB.SetConnMaxLifetime(time.Duration(customDB.cfg.Pool.MaxLifeTime) * time.Second)
		}
	}
}

func (customDB *DB) validateCfg() error {
	customDB.cfgSeqValidatorChains = append(
		customDB.cfgSeqValidatorChains,
		customDB.cfgValidatorNil, customDB.cfgValidatorTypeEmpty, customDB.cfgValidatorDSNEmpty,
	)
	for _, cfgValidator := range customDB.cfgSeqValidatorChains {
		if err := cfgValidator(); err != nil {
			return err
		}
	}
	return nil
}

func (customDB *DB) cfgValidatorNil() error {
	if customDB.cfg == nil {
		return ErrorDBCfgNotFound
	}
	return nil
}

func (customDB *DB) cfgValidatorTypeEmpty() error {
	if customDB.cfg.Type == "" {
		return ErrorDBCfgTypeNotFound
	}
	return nil
}

func (customDB *DB) cfgValidatorDSNEmpty() error {
	if customDB.cfg.Dsn == "" {
		return ErrorDBCfgDSNNotFound
	}
	return nil
}
