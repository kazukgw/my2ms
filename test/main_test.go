package main

import (
	"math/rand"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
)

func TestSameData(t *testing.T) {
	ctx := SetUp()

	mydb := ctx.Value("mysql").(*gorm.DB)
	msdb := ctx.Value("mssql").(*gorm.DB)

	users := []*User{}
	mydb.Find(&users)

	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		expect := users[rand.Intn(len(users))]

		log.Info("find user on mssql")
		us := []*User{}
		msdb.Where("user_id = ?", expect.UserID).Find(&us)

		if len(us) < 1 {
			t.Error(spew.Sprintf(
				"user not found. user: %s",
				expect,
			))
		}

		u := us[0]
		if expect.Code != u.Code ||
			expect.Name != u.Name ||
			expect.Email != u.Email {
			t.Errorf(
				"expect: %s, actual: %s",
				spew.Sprint(expect),
				spew.Sprint(u),
			)
		}
	}
}
