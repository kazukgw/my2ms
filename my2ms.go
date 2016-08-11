package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/siddontang/go-mysql/canal"
	"github.com/spf13/viper"
)

func main() {
	configPath := flag.String("confdir", ".", "conf directory path")
	flag.Parse()

	conf := InitConfig(*configPath)

	log := InitLogger(conf)

	log.Info("start my2ms")

	log.Info("initialize a canal")
	c := InitCanal(conf, log)

	log.Info("create new handler")
	h := NewMy2MSHandler(conf, log)
	c.RegRowsEventHandler(h)
	if err := c.Start(); err != nil {
		panic(err.Error())
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(
		sigchan,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	for {
		select {
		case <-sigchan:
			break
		}
	}

	log.Println("quit my2ms")
}

func InitConfig(configPath string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigType("yaml")
	conf.SetConfigName("my2ms")
	conf.AddConfigPath(configPath)
	if err := conf.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	return conf
}

func InitLogger(conf *viper.Viper) *logrus.Logger {
	logger := logrus.New()
	level := conf.GetString("log.level")
	logout := conf.GetString("log.out")
	var out io.Writer
	if logout == "" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.OpenFile(logout, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
	logger.Out = out
	switch strings.ToLower(level) {
	case "debug":
		logger.Level = logrus.DebugLevel
	case "info":
		logger.Level = logrus.InfoLevel
	default:
		logger.Level = logrus.InfoLevel
	}
	logger.Level = logrus.DebugLevel
	return logger
}

func InitCanal(conf *viper.Viper, logger *logrus.Logger) *canal.Canal {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = conf.GetString("canal.host")
	canalCfg.User = conf.GetString("canal.user")
	canalCfg.Password = conf.GetString("canal.password")
	canalCfg.Dump.TableDB = conf.GetString("canal.table_db")
	canalCfg.Dump.Tables = conf.GetStringSlice("canal.tables")
	canalCfg.Dump.ExecutionPath = conf.GetString("canal.dump_execution_path")

	logger.Debug(spew.Sprintf("canal config:", canalCfg))
	c, err := canal.NewCanal(canalCfg)
	if err != nil {
		panic(err.Error())
	}
	if err := c.CheckBinlogRowImage("FULL"); err != nil {
		panic(err.Error())
	}
	return c
}

func NewMy2MSHandler(conf *viper.Viper, logger *logrus.Logger) *My2MSHandler {
	h := &My2MSHandler{}
	dns := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s",
		conf.GetString("mssql.server"),
		conf.GetString("mssql.database"),
		conf.GetString("mssql.user"),
		conf.GetString("mssql.password"),
	)

	logger.Debug(spew.Sprintf("mssql dns:", dns))
	db, err := sql.Open("mssql", dns)
	if err != nil {
		panic(err.Error())
	}
	logger.Debug(spew.Sprintf("mssql connection:", db))
	h.MSDB = db

	sb := &SqlBuilder{}
	mmaps := NewMigrationMapsWithMap(conf.GetStringMapString("migration_map"))
	logger.Debug(spew.Sprintf("my2ms table maps:", mmaps))
	sb.MigrationMaps = mmaps
	h.SqlBuilder = sb
	h.Logger = logger

	return h
}

type My2MSHandler struct {
	MSDB *sql.DB
	*SqlBuilder
	Logger *logrus.Logger
}

func (My2MSHandler) String() string {
	return "My2MSHandler"
}

func (h *My2MSHandler) Do(e *canal.RowsEvent) error {
	sqlsets := h.SqlBuilder.BuildSql(e)
	for _, sqlset := range *sqlsets {
		if sqlset.Error != nil {
			return sqlset.Error
		}
		_, err := h.MSDB.Exec(sqlset.Sql, sqlset.Args...)
		if err != nil {
			return err
		}
	}
	return nil
}
