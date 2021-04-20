package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"short-url-server/comm"
	"short-url-server/common"
	"short-url-server/repo"
	"short-url-server/repo/model"
)

// UrlHandler godoc
// @Summary Read URLs
// @Description get all URLs from DB
// @Produce json
// @Success 200 {object} comm.AllUrlResponse
// @Failure 500 {object} comm.FailureResponse
// @Router /urls [get]
func UrlHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if urls, err := db.FindAllUrls(); err != nil {
			c.JSON(http.StatusInternalServerError, comm.FailureResponse{
				Status:  common.Failed,
				Message: "retrieving urls failed",
			})
		} else {
			c.JSON(http.StatusOK, comm.AllUrlResponse{
				Status: common.Succeed,
				Urls:   convertToUrlEntities(urls),
			})
		}
	}
}

// AddUrlHandler godoc
// @Summary Create shorten URL
// @Description Create new URL mapping
// @Accept json
// @Produce json
// @Param req body comm.CreateUrlRequest true "Create short URL"
// @Success 202 {object} comm.CreateUrlResponse
// @Failure 406 {object} comm.FailureResponse
// @Failure 500 {object} comm.FailureResponse
// @Router /urls [post]
func AddUrlHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req comm.CreateUrlRequest
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("while read: %v", err),
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("while unmarshal: %v", err),
			})
		} else if id, err := insertUrlUntilSuccess(db, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("while insert to db: %v", err),
			})
		} else {
			c.JSON(http.StatusAccepted, comm.CreateUrlResponse{
				Status: common.Succeed,
				Id:     *id,
			})
		}
	}
}

// SingleUrlHandler godoc
// @Summary URL info
// @Description Get Single URL Information
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Success 200 {object} comm.GetUrlResponse
// @Failure 404 {object} comm.FailureResponse
// @Router /urls/{id} [get]
func SingleUrlHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if url, err := db.FindUrlById(id); err != nil {
			c.JSON(http.StatusNotFound, comm.FailureResponse{
				Status:  common.Failed,
				Message: "id not found",
			})
		} else {
			c.JSON(http.StatusOK, comm.GetUrlResponse{
				Status: common.Succeed,
				Url:    convertToUrlEntity(*url),
			})
		}
	}
}

// SingleUrlUpdateHandler godoc
// @Summary update URL entry
// @Description Update single URL entry
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Param request body comm.UpdateUrlRequest true "URL to be updated"
// @Success 200 {object} comm.UpdateUrlResponse
// @Failure 406 {object} comm.FailureResponse
// @Failure 500 {object} comm.FailureResponse
// @Router /urls/{id} [patch]
func SingleUrlUpdateHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req comm.UpdateUrlRequest
		if bytes, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("error when parsing request body: %v", err),
			})
		} else if err := json.Unmarshal(bytes, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("error while unmarshaling request body: %v", err),
			})
		} else if url, err := db.UpdateUrl(id, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("error while updating database: %v", err),
			})
		} else {
			c.JSON(http.StatusOK, comm.UpdateUrlResponse{
				Status:  common.Succeed,
				Updated: url > 0,
			})
		}
	}
}

// SingleUrlDeleteHandler godoc
// @Summary delete URL entry
// @Description delete single URL entry
// @Accept json
// @Produce json
// @Param id path string true "URL ID"
// @Success 202 {object} comm.SuccessResponse
// @Failure 406 {object} comm.FailureResponse
// @Failure 500 {object} comm.FailureResponse
// @Router /urls/{id} [delete]
func SingleUrlDeleteHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id = c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotAcceptable, comm.FailureResponse{
				Status:  common.Failed,
				Message: "path should be not empty",
			})
		} else if err := db.DeleteUrl(id); err != nil {
			c.JSON(http.StatusInternalServerError, comm.FailureResponse{
				Status:  common.Failed,
				Message: fmt.Sprintf("Error while delete from DB: %+v", err),
			})
		} else {
			c.JSON(http.StatusAccepted, comm.SuccessResponse{
				Status:  common.Succeed,
				Message: "deleted successfully!",
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

func convertToUrlEntity(src model.Url) comm.UrlEntity {
	return comm.UrlEntity{
		Id:        src.ShortenId,
		Url:       src.Url,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func convertToUrlEntities(srcs []model.Url) []comm.UrlEntity {
	ret := make([]comm.UrlEntity, len(srcs))
	for idx, src := range srcs {
		ret[idx] = convertToUrlEntity(src)
	}
	return ret
}
