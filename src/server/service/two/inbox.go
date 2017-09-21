package two


import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/two"
	"github.com/jinzhu/gorm"
)

type InboxService struct {

}

func (self *InboxService)GetByID(id int) (*two.Inbox, error) {
	inbox := two.Inbox{}
	err := common.Soda2DB_R.Where(id).Find(&inbox).Error
	if err != nil {
		return nil, err
	}
	return &inbox, nil
}

//CountConsultation 获取询问商品人数
func (self *InboxService)CountConsultation(topicID int, userID int) (int, error) {
	count := 0
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if topicID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("target_id = ?", topicID)
		})
	}

	if userID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("receiver_id = ?", userID)
		})
	}
	err := common.Soda2DB_R.Table("inbox").Scopes(scopes...).Where("is_official = 0 and type = 2").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil

}

