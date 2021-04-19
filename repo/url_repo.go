package repo

import (
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

func (r *DDB) DeleteUrl(id string) error {
	if result := r.DB().Delete(&model.Url{}, "shorten_id = ?", id); result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
