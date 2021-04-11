package repo

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

var db *DDB

func init() {
	var logger = logrus.New()
	db = New(logger)
}

func TestInsertUrl(t *testing.T) {
	var err error
	fmt.Println("insert!")
	err = db.InsertUrl("abc", "def")
	if err != nil {
		t.Fail()
	}
	fmt.Println("find!")
	url, err := db.FindUrlById("abc")
	if err != nil || url.Url != "def" {
		t.Fail()
	}
	fmt.Println("about to delete!")
	err = db.DeleteUrl("abc")
	if err != nil {
		t.Fail()
	}
	fmt.Println("find again should fail!")
	_, w := db.FindUrlById("abc")
	if w != nil {
		t.Fail()
	}
}
