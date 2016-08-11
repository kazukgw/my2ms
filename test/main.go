package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/manveru/faker"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const MYSQL_SERVER = "192.168.99.100:3306"

func main() {
	action := flag.String("action", "", "")
	flag.Parse()

	switch *action {
	case "reset":
		Reset()
	case "insert":
		Insert()
	case "update":
		Update()
	case "delete":
		Delete()
	}
}

func SetUp() context.Context {
	conf := viper.New()
	conf.SetConfigFile("mysql/app/my2ms.yml")
	if err := conf.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "conf", conf)

	mydb := MustOpenMySQL(conf, MYSQL_SERVER)
	ctx = context.WithValue(ctx, "mysql", mydb)

	msdb := MustOpenMSSQL(conf)
	ctx = context.WithValue(ctx, "mssql", msdb)

	return ctx
}

func Reset() {
	ctx := SetUp()

	mydb := ctx.Value("mysql").(*gorm.DB)
	mydb.DropTableIfExists(&User{})
	mydb.CreateTable(&User{})

	msdb := ctx.Value("mssql").(*gorm.DB)
	msdb.DropTableIfExists(&User{})
	msdb.CreateTable(&User{})
}

func Insert() {
	ctx := SetUp()

	mydb := ctx.Value("mysql").(*gorm.DB)
	for i := 0; i < 1000; i++ {
		u := MustGenRandomUser()
		mydb.Create(u)
	}
}

func Update() {
	ctx := SetUp()

	fake, err := faker.New("en")
	if err != nil {
		panic(err.Error())
	}

	mydb := ctx.Value("mysql").(*gorm.DB)
	users := []*User{}
	mydb.Find(&users)
	for _, u := range users {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(3) == 2 {
			u.Email = fake.Email()
			u.Code = fake.Rand.Intn(1000)
			u.Name = fake.UserName()
			mydb.Save(u)
		}
	}
}

func Delete() {
	ctx := SetUp()

	mydb := ctx.Value("mysql").(*gorm.DB)
	users := []*User{}
	mydb.Find(&users)

	cnt := 0
	for _, u := range users {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(3) == 2 {
			mydb.Delete(u)
			cnt++
		}
	}
	log.Debug("delete %i records", cnt)
}

func MustOpenMySQL(conf *viper.Viper, host string) *gorm.DB {
	// host := conf.GetString("canal.host")
	user := conf.GetString("canal.user")
	password := conf.GetString("canal.password")
	database := conf.GetString("canal.table_db")
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		user,
		password,
		host,
		database,
	)
	log.Debug(spew.Sprintf("mysql dns:", dns))
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	return db
}

func MustOpenMSSQL(conf *viper.Viper) *gorm.DB {
	dns := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s",
		conf.GetString("mssql.server"),
		conf.GetString("mssql.database"),
		conf.GetString("mssql.user"),
		conf.GetString("mssql.password"),
	)
	log.Debug(spew.Sprintf("mssql dns:", dns))
	db, err := gorm.Open("mssql", dns)
	if err != nil {
		panic(err)
	}
	return db
}

func MustGenRandomUser() *User {
	fake, err := faker.New("en")
	if err != nil {
		panic(err.Error())
	}
	u := &User{}
	u.UserID = uuid.New().String()
	u.Name = fake.UserName()
	u.Code = fake.Rand.Intn(11)
	u.Email = fake.Email()
	return u
}
