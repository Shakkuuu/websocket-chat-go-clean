package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error

type Postgres struct {
	Db *gorm.DB
}

func New(host, user, password, database, dbport string) (*Postgres, error) {
	pg := &Postgres{}

	CONNECT := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, database, dbport)

	// 接続できるまで一定回数リトライ
	count := 0
	pg.Db, err = gorm.Open(postgres.Open(CONNECT), &gorm.Config{})
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 { // countが180になるまでリトライ
				fmt.Println("")
				log.Printf("db Init error: %v\n", err)
				panic(err)
			}
			pg.Db, err = gorm.Open(postgres.Open(CONNECT), &gorm.Config{})
		}
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if sqlDB, err := p.Db.DB(); err != nil {
		log.Printf("db Close error: %v\n", err)
		panic(err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("db Close error: %v\n", err)
			panic(err)
		}
	}
}
