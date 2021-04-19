package repo

import (
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"short-url-server/repo/model"
)

type DDB struct {
	db  *gorm.DB
}

const DSN = "shortener:shortener@tcp(127.0.0.1:3306)/shortener?charset=utf8mb4&parseTime=True&loc=Local"

func New(log *logrus.Logger) *DDB {
	logrus.SetLevel(logrus.DebugLevel)
	if db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		Logger: gorm_logrus.New(),
	}); err != nil {
		log.Fatalf("Failed to open server: %v", err)
		panic("Can't reach here")
	} else if err := db.AutoMigrate(&model.Url{}); err != nil {
		log.Fatalf("Failed to migrate table: %v", err)
		panic("Can't reach here")
	} else {
		return &DDB{
			db:  db,
		}
	}
}

func (r *DDB) DB() *gorm.DB {
	return r.db
}
