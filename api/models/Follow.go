package models

import (
	"errors"
	_ "time"

	"github.com/jinzhu/gorm"
)

type Follow struct {
	ID         uint64 `gorm:"primary_key;auto_increment" json:"id"`
	FollowerID uint32 `gorm:"not null" json:"follower_id"`
	FolloweeID uint32 `gorm:"not null" json:"followee_id"`
	User       User   `json:user`
}

func (f *Follow) Prepare() {
	f.ID = 0
	f.FollowerID = 0
	f.FolloweeID = 0
	f.User = User{}
}

func (f *Follow) CreateFollow(db *gorm.DB, followeeID int32) (*Follow, error) {
	var err error
	err = db.Debug().Model(&Follow{}).Create(&f).Error
	if err != nil {
		return &Follow{}, err
	}
	if f.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", f.FollowerID).Take(&f.User).Error
		if err != nil {
			return &Follow{}, err
		}
	}
	return f, nil
}

func (f *Follow) FindAllFollows(db *gorm.DB) (*[]Follow, error) {
	var err error
	follows := []Follow{}
	err = db.Debug().Model(&Follow{}).Limit(100).Find(&follows).Error
	if err != nil {
		return &[]Follow{}, err
	}
	if len(follows) > 0 {
		for i, _ := range follows {
			err := db.Debug().Model(&User{}).Where("id = ?", follows[i].FollowerID).Take(&follows[i].User).Error
			if err != nil {
				return &[]Follow{}, err
			}
		}
	}
	return &follows, nil
}

func (f *Follow) FindFollowByID(db *gorm.DB, fid uint64) (*Follow, error) {
	var err error
	err = db.Debug().Model(&Follow{}).Where("id = ?", fid).Take(&f).Error
	if err != nil {
		return &Follow{}, err
	}
	if f.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", f.FollowerID).Take(&f.User).Error
		if err != nil {
			return &Follow{}, err
		}
	}
	return f, nil
}

func (f *Follow) DeleteAFollow(db *gorm.DB, fid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Follow{}).Where("id = ? and follower_id = ?", fid, uid).Take(&Follow{}).Delete(&Follow{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Follow not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
