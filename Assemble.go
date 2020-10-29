package models

import (
	"strings"
	"time"
	"xiaofendian/common"


	哈哈哈哈哈哈
	第2次合并

	master
)

type Assemble struct {
	Id         		  uint64        `gorm:"column:id;primary_key" json:"id"`
	ActId             uint64        `gorm:"column:act_id;type:bigint(20)" json:"act_id" form:"act_id"`//活动id
	Status        	  int       	`gorm:"cloumn:status;type:int(10)" json:"status" form:"status"` //0:关闭 1:正常
	Uid           	  uint64        `gorm:"column:uid;type:bigint(20)" json:"uid" form:"uid"`
	ShopId        	  uint64    	`gorm:"column:shop_id;type:bigint(20)" json:"shop_id" form:"shop_id"`
	Title      		  string 		`gorm:"column:title;type:varchar(100)" json:"title" form:"title"`
	Banner     		  string        `gorm:"column:banner;type:varchar(500)" json:"banner" form:"banner"`
	NeedPhone         int           `gorm:"column:need_phone;type:int(10)" json:"need_phone" form:"need_phone"`
	GoodsNums         int           `gorm:"column:goods_nums;type:int(10)" json:"goods_nums" form:"goods_nums"`//商品数量
	GoodsPrice        uint64         `gorm:"column:goods_price;type:bigint(20)" json:"goods_price" form:"goods_price"`//商品原价
	GroupPrice        uint64 		`gorm:"column:group_price;type:bigint(20)" json:"group_price" form:"group_price"`//拼团价
	FullPay           int        	`gorm:"column:full_pay;type:int(10)" json:"full_pay" form:"full_pay"`//全额付款
	PrepayAmount      uint64         `gorm:"column:prepay_amount;type:bigint(20);DEFAULT:0"  json:"prepay_amount" form:"prepay_amount"`//预付款
	GroupNums         int           `gorm:"column:group_nums;type:int(10)" json:"group_nums" form:"group_nums"`//成团人数
	IsSmart           int           `gorm:"column:is_smart;type:int(10)" json:"is_smart" form:"is_smart"`//智能防刷
	Describe     	  string        `gorm:"column:describe;type:varchar(1500)" json:"describe" form:"describe"`
	MusicId           uint64        `gorm:"column:music_id;type:int(10)" json:"music_id" form:"music_id"`
	Rules     		  string        `gorm:"column:rules;type:varchar(2000)" json:"rules" form:"rules"`
	StartTime         time.Time 	`gorm:"type:datetime" json:"start_time" time_format:"2006-01-02 15:04:05"  form:"start_time"`
	EndTime           time.Time 	`gorm:"type:datetime" json:"end_time" time_format:"2006-01-02 15:04:05" form:"end_time"`
	Base
	Activity 		  *ActivityList  `gorm:"ForeignKey:ActId;AssociationForeignKey:Id" json:"activity,omitempty" form:"activity"`
	Music 		      *Music         `gorm:"ForeignKey:MusicId;AssociationForeignKey:Id" json:"music,omitempty" form:"music"`
	MusicName         string    	 `gorm:"-" json:"musci_name,omitempty" form:"musci_name"`
	Multimedias 	  []Multimedia   `gorm:"ForeignKey:ActId;AssociationForeignKey:ActId" json:"multimedias" form:"multimedias"`
	Shop               *Shops     	 `gorm:"ForeignKey:ShopId;AssociationForeignKey:Id" json:"shop,omitempty"`
	JoinNums		  int			 `gorm:"-" json:"join_nums" form:"join_nums"`

}

func test() {
	var_dump("nihao,world");
}

func (Assemble) TableName() string  {
	return "assemble"
}

/**
商家创建拼团活动
 */
func CreateAssembleActivity(uid uint64,assemble Assemble) (id uint64,err error) {
	local, _ := time.LoadLocation("Local")
	start_time, _ := time.ParseInLocation(common.DATE_FORMAT, assemble.StartTime.Format(common.DATE_FORMAT), local)
	end_time, _ := time.ParseInLocation(common.DATE_FORMAT, assemble.EndTime.Format(common.DATE_FORMAT), local)
	o := Gorms()["default"]
	tx := o.Begin()
	var activity ActivityList
	activity.Uid = uid
	activity.ShopId = assemble.ShopId
	activity.Type = 7 //拼团活动
	err = tx.Create(&activity).Error
	if err != nil {
		tx.Rollback()
		return
	}
	var assem Assemble
	assem.ActId = activity.Id
	assem.Title = assemble.Title
	if assemble.Banner != "" {
		assem.Banner = common.Img_Url + assemble.Banner
	}
	assem.Status = 1
	assem.NeedPhone = assemble.NeedPhone
	assem.FullPay = assemble.FullPay
	assem.PrepayAmount = assemble.PrepayAmount
	assem.GroupNums = assemble.GroupNums
	assem.GoodsNums = assemble.GoodsNums
	assem.Uid = uid
	assem.GoodsPrice = assemble.GoodsPrice
	assem.GroupPrice = assemble.GroupPrice
	assem.StartTime = start_time
	assem.EndTime = end_time
	assem.Describe = assemble.Describe
	assem.MusicId = assemble.MusicId
	assem.Rules = assemble.Rules
	assem.IsSmart = assemble.IsSmart
	err = tx.Create(&assem).Error
	if err != nil {
		tx.Rollback()
		return
	}
	if len(assemble.Multimedias) >0 {
		for _, img := range assemble.Multimedias {
			newImg := Multimedia{}
			newImg.ActId = activity.Id
			if img.Url != "" {
				newImg.Url = common.Img_Url + img.Url
			}
			newImg.MediaType = 1
			err = tx.Create(&newImg).Error
			if err != nil {
				tx.Rollback()
				return
			}
		}
	}
	id = activity.Id
	tx.Commit()
	return
}

func UpdateAssembleActivity(assemble Assemble) (err error) {
	save := make(map[string]interface{})
	local, _ := time.LoadLocation("Local")
	start_time, _ := time.ParseInLocation(common.DATE_FORMAT, assemble.StartTime.Format(common.DATE_FORMAT), local)
	end_time, _ := time.ParseInLocation(common.DATE_FORMAT, assemble.EndTime.Format(common.DATE_FORMAT), local)
	o := Gorms()["default"]
	tx := o.Begin()

	save["title"] = assemble.Title
	if  strings.Contains(assemble.Banner,"https://") {
		save["banner"] = assemble.Banner
	}else {
		save["banner"]  = common.Img_Url + assemble.Banner
	}
	save["need_phone"] = assemble.NeedPhone
	save["goods_nums"] = assemble.GoodsNums
	save["goods_price"] = assemble.GoodsPrice
	save["group_price"] = assemble.GroupPrice
	save["full_pay"] = assemble.FullPay
	save["prepay_amount"] = assemble.PrepayAmount
	save["group_nums"] = assemble.GroupNums
	save["is_smart"] = assemble.IsSmart
	save["describe"] = assemble.Describe
	save["music_id"] = assemble.MusicId
	save["rules"] = assemble.Rules
	save["start_time"] = start_time
	save["end_time"] = end_time
	err = tx.Model(&Assemble{}).Where("act_id = ?",assemble.ActId).Updates(save).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Where("act_id = ?",assemble.ActId).Delete(Multimedia{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	if len(assemble.Multimedias) > 0 {
		for _, img := range assemble.Multimedias {
			newImg := Multimedia{}
			newImg.ActId = assemble.ActId
			if img.Url != "" {
				if  strings.Contains(img.Url,"https://") {
					newImg.Url  = img.Url
				}else {
					newImg.Url = common.Img_Url + img.Url
				}
			}
			newImg.MediaType = 1
			err = tx.Create(&newImg).Error
			if err != nil {
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
	return
}