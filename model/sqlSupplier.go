package model

import (
	"luckgo/log"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"sync/atomic"
	"time"
)

const (
	EXIT_DB_OPEN = 101
	EXIT_PING    = 102
)

const (
	MAX_DB_CONN_LIFETIME = 60
	DB_PING_ATTEMPTS     = 18
	DB_PING_TIMEOUT_SECS = 10
	DB_OPEN_ATTEMPTS     = 18
	DB_OPEN_TIMEOUT_SECS = 10	
)

type SqlSupplier struct {
	master         *gorm.DB
	replicas       []*gorm.DB
	searchReplicas []*gorm.DB
	rrCounter      int64
	srCounter      int64
}

var sqlSupplier *SqlSupplier

func NewSqlSupplier() *SqlSupplier {
	sqlSupplier = &SqlSupplier{
		rrCounter: 0,
		srCounter: 0,
	}

	sqlSupplier.InitConnection()
	//InitTable(sqlSupplier)
	return sqlSupplier
}

func (ss *SqlSupplier) InitConnection() {
	ss.master = setupConnection("master", *Cfg.SqlSettings.DriverName, *Cfg.SqlSettings.DataSource,
		*Cfg.SqlSettings.MaxIdleConns, *Cfg.SqlSettings.MaxOpenConns, *Cfg.SqlSettings.Trace)
	if len(Cfg.SqlSettings.DataSourceReplicas) == 0 {
		ss.replicas = make([]*gorm.DB, 1)
		ss.replicas[0] = ss.master
	} else {
		ss.replicas = make([]*gorm.DB, len(Cfg.SqlSettings.DataSourceReplicas))
		for i, replica := range Cfg.SqlSettings.DataSourceReplicas {
			ss.replicas[i] = setupConnection(fmt.Sprintf("replica-%v", i),
				*Cfg.SqlSettings.DriverName, replica, *Cfg.SqlSettings.MaxIdleConns, *Cfg.SqlSettings.MaxOpenConns, *Cfg.SqlSettings.Trace)
		}
	}

	if len(Cfg.SqlSettings.DataSourceSearchReplicas) == 0 {
		ss.searchReplicas = ss.replicas
	} else {
		ss.searchReplicas = make([]*gorm.DB, len(Cfg.SqlSettings.DataSourceSearchReplicas))
		for i, replica := range Cfg.SqlSettings.DataSourceSearchReplicas {
			ss.searchReplicas[i] = setupConnection(fmt.Sprintf("search-replica-%v", i),
				*Cfg.SqlSettings.DriverName, replica, *Cfg.SqlSettings.MaxIdleConns, *Cfg.SqlSettings.MaxOpenConns, *Cfg.SqlSettings.Trace)
		}
	}
	return
}

func setupConnection(con_type string, driver string, dataSource string, maxIdle int, maxOpen int, trace bool) *gorm.DB {
	log.Info("driver:" + con_type)
	log.Info("open database:" + dataSource)
	var db *gorm.DB
	var err error
	for i := 0; i < DB_OPEN_ATTEMPTS; i++ {
		db, err = gorm.Open(driver, dataSource)
		if err == nil {
			log.Info("open database successfully")
			break
		} else {
			if i == DB_OPEN_ATTEMPTS - 1 {
				log.Critical(err.Error())
				time.Sleep(time.Second)
				os.Exit(EXIT_DB_OPEN)
				return nil 
			} else {
				log.Error(fmt.Sprintf("Failed to Open DB retrying  in %v seconds err = %v", DB_OPEN_TIMEOUT_SECS, err))
				time.Sleep(DB_OPEN_TIMEOUT_SECS * time.Second)
			}
		}
	}

	for i := 0; i < DB_PING_ATTEMPTS; i++ {
		log.Info("Pinging SQL " + con_type + " database")
		if err := db.DB().Ping(); err == nil {
			break
		} else {
			if i == DB_PING_ATTEMPTS-1 {
				log.Critical("Failed to ping DB, server will exit err =" + err.Error())
				time.Sleep(time.Second)
				os.Exit(EXIT_PING)
				return nil
			} else {
				log.Error(fmt.Sprintf("Failed to ping DB retrying  in %v seconds err = %v", DB_PING_TIMEOUT_SECS, err))
				time.Sleep(DB_PING_TIMEOUT_SECS * time.Second)
			}
		}
	}
	//设置数据库的空闲连接和最大打开连接
	db.DB().SetMaxIdleConns(*Cfg.SqlSettings.MaxIdleConns)
	db.DB().SetMaxOpenConns(*Cfg.SqlSettings.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(MAX_DB_CONN_LIFETIME) * time.Minute)
	db.LogMode(*Cfg.SqlSettings.Trace)
	/*
		connectionTimeout := time.Duration(*utils.Cfg.Cfg.SqlSettings.QueryTimeout) * time.Second

		if trace {
			dbmap.TraceOn("", sqltrace.New(os.Stdout, "sql-trace:", sqltrace.Lmicroseconds))
		}
	*/
	return db
}
func (ss *SqlSupplier) GetMaster() *gorm.DB {
	return ss.master
}

func (ss *SqlSupplier) GetSearchReplica() *gorm.DB {
	rrNum := atomic.AddInt64(&ss.srCounter, 1) % int64(len(ss.searchReplicas))
	return ss.searchReplicas[rrNum]
}

func (ss *SqlSupplier) GetReplica() *gorm.DB {
	rrNum := atomic.AddInt64(&ss.rrCounter, 1) % int64(len(ss.replicas))
	return ss.replicas[rrNum]
}

func (ss *SqlSupplier) TotalMasterDbConnections() int {
	return ss.GetMaster().DB().Stats().OpenConnections
}

func (ss *SqlSupplier) TotalReadDbConnections() int {

	if len(Cfg.SqlSettings.DataSourceReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.replicas {
		count = count + db.DB().Stats().OpenConnections
	}

	return count
}

func (ss *SqlSupplier) TotalSearchDbConnections() int {
	if len(Cfg.SqlSettings.DataSourceSearchReplicas) == 0 {
		return 0
	}

	count := 0
	for _, db := range ss.searchReplicas {
		count = count + db.DB().Stats().OpenConnections
	}

	return count
}

func (ss *SqlSupplier) Close() {
	log.Info("Closing SqlStore")
	ss.master.DB().Close()
	for _, replica := range ss.replicas {
		replica.DB().Close()
	}
}
