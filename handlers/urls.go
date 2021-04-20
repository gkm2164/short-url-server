package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"short-url-server/common"
	"short-url-server/repo"
)

func UrlHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ret, err := db.FindAllUrls(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  common.Failed,
				"message": "retrieving urls failed",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": common.Succeed,
				"urls":   ret,
			})
		}
	}
}


func AddUrlHandler(db *repo.DDB) gin.HandlerFunc {
	type NewAddrRequest struct {
		Url string `json:"url"`
	}

	return func(c *gin.Context) {
		var req NewAddrRequest
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status": common.Failed,
				"message": fmt.Sprintf("while read: %v", err),
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status": common.Failed,
				"message": fmt.Sprintf("while unmarshal: %v", err),
			})
		} else if id, err := insertUrlUntilSuccess(db, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": common.Failed,
				"message": fmt.Sprintf("while insert to db: %v", err),
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"status": common.Succeed,
				"id": id,
			})
		}
	}
}

// SingleUrlHandler route == /urls/:id
func SingleUrlHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if url, err := db.FindUrlById(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  common.Failed,
				"message": "id not found",
			})
		} else {
			c.JSON(http.StatusFound, gin.H{
				"status": common.Succeed,
				"url":    url,
			})
		}
	}
}

// SingleUrlUpdateHandler route == /urls/:id
func SingleUrlUpdateHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		type Request struct {
			Url string `json:"url"`
		}
		var req Request
		if bytes, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  common.Failed,
				"message": fmt.Errorf("error when parsing request body: %v", err),
			})
		} else if err := json.Unmarshal(bytes, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  common.Failed,
				"message": fmt.Errorf("error while unmarshaling request body: %v", err),
			})
		} else if url, err := db.UpdateUrl(id, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  common.Failed,
				"message": fmt.Errorf("error while updating database: %v", err),
			})
		} else {
			c.JSON(http.StatusFound, gin.H{
				"status": common.Succeed,
				"url":    url,
			})
		}
	}
}

func SingleUrlDeleteHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id = c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  common.Failed,
				"message": "path should be not empty",
			})
		} else if err := db.DeleteUrl(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  common.Failed,
				"message": err,
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"status":  common.Succeed,
				"message": "deleted successfully!",
			})
		}
	}
}

const UrlSize = 11

func insertUrlUntilSuccess(db *repo.DDB, url string) (*string, error) {
	var success = false
	for !success {
		id := randomString()
		if _, err := db.InsertUrl(id, url); err != nil {
			return nil, fmt.Errorf("errors when trying to put item: %v", err)
		} else {
			success = true
			return &id, nil
		}
	}

	return nil, errors.New("should not reach here")
}

func randomString() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	retStr := make([]uint8, UrlSize)
	for i := 0; i < UrlSize; i++ {
		retStr[i] = letters[rand.Uint32()%uint32(len(letters))]
	}
	return string(retStr)
}
