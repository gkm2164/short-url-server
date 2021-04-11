package repo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"short-url/repo/model"
	"time"
)

var TableName = aws.String("ShortUrls")

func (r *DDB) FindUrlById(id string) (*model.Url, error) {
	var url model.Url
	if result, err := r.DB().GetItem(&dynamodb.GetItemInput{
		TableName: TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Item, &url); err != nil {
		return nil, err
	} else {
		return &url, nil
	}

}

func (r *DDB) InsertUrl(id string, url string) error {
	if av, err := dynamodbattribute.MarshalMap(model.Url{
		Id:        id,
		Url:       url,
		CreatedAt: time.Now(),
	}); err != nil {
		r.Log().Errorf("error %v", err)
		return err
	} else if _, err := r.DB().PutItem(&dynamodb.PutItemInput{
		TableName: TableName,
		Item: av,
	}); err != nil {
		return err
	}

	return nil
}

func (r *DDB) DeleteUrl(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: TableName,
	}

	_, err := r.DB().DeleteItem(input)
	if err != nil {
		r.Log().Errorf("failed to delete url: %s", err)
		return err
	}

	return nil
}
