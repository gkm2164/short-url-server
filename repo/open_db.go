package repo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

type DDB struct {
	db *dynamodb.DynamoDB
	log *logrus.Logger
}

func New(log *logrus.Logger) *DDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &DDB{
		db: dynamodb.New(sess),
		log: log,
	}
}

func (r *DDB) DB() *dynamodb.DynamoDB {
	return r.db
}

func (r *DDB) Log() *logrus.Logger {
	return r.log
}
