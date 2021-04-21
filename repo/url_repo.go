package repo

import (
	"errors"
	"fmt"
	"short-url-server/repo/model"
)

func (r *DDB) FindUrlById(id string) (*model.Url, error) {
	var url model.Url
	if tx := r.DB().First(&url, "shorten_id = ?", id); tx.Error != nil {
		return nil, tx.Error
	} else {
		return &url, nil
	}
}

func (r *DDB) IncAccessCount(id string) error {
	if tx := r.DB().Exec("UPDATE urls SET access_count = access_count + 1 WHERE shorten_id = ?", id); tx == nil {
		return errors.New("error while update")
	} else if tx.Error != nil {
		return fmt.Errorf("error from db: %v", tx.Error)
	} else {
		return nil
	}
}

func (r *DDB) FindAllUrls() ([]model.Url, error) {
	var urls []model.Url
	if tx := r.DB().Find(&urls); tx.Error != nil {
		return nil, tx.Error
	} else {
		return urls, nil
	}
}

func (r *DDB) InsertUrl(id string, url string) (uint, error) {
	urlEntity := model.Url{
		ShortenId: id,
		Url:       url,
	}

	if result := r.DB().Create(&urlEntity); urlEntity.ID == 0 {
		return urlEntity.ID, result.Error
	} else {
		return 0, nil
	}
}

func (r *DDB) UpdateUrl(id string, url string) (int64, error) {
	if tx := r.DB().Model(&model.Url{}).
		Where("shorten_id = ?", id).
		Update("url", url); tx == nil {
		return 0, errors.New("transaction returns null")
	} else if tx.Error != nil {
		return 0, fmt.Errorf("error from DB: %v", tx.Error)
	} else {
		return tx.RowsAffected, nil
	}
}

func (r *DDB) DeleteUrl(id string) error {
	if result := r.DB().Delete(&model.Url{}, "shorten_id = ?", id); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
