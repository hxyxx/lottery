package dao

import (
	"github.com/go-xorm/xorm"
	"hxyxx/lottery/models"
	"log"
)

type GiftDao struct {
	engine *xorm.Engine
}

func NewGiftDao(engine *xorm.Engine) *GiftDao{
	return &GiftDao{engine: engine}
}

func (d *GiftDao) Get(id int) *models.LtGift{
	//构建一个model对象
	 data := &models.LtGift{Id :id}
	 //查询
	 ok,err := d.engine.Get(data)
	 if ok && err !=nil{
	 	return data
	 }else {
	 	data.Id = 0
	 	return data
	 }
}
func (d *GiftDao) GetAll() []models.LtGift{
	datalist := make([]models.LtGift,0)
	err := d.engine.
		Asc("sys_status").
		Asc("displayorder").Find(&datalist)
	if err != nil {
		log.Println("gift_dao.GetAll error = ",err)
	}
	return datalist
}
func (d *GiftDao) CountAll() int64 {
	num,err := d.engine.Count(&models.LtGift{})
	if err !=nil{
		log.Println("gift_dao.CountAll error : ",err)
		return 0
	}
	return num
}
func (d *GiftDao) Delete(id int) error{
	data := &models.LtGift{Id:id}
	_,err := d.engine.Delete(data)
	return err
}
func (d *GiftDao) update(data *models.LtGift,columns []string) error{
	//MustCols方法指定对象属性即使为空也进行更新，指定字段
	_,err := d.engine.ID(data.Id).MustCols(columns...).Update(data)
	return err
}
func (d *GiftDao) create(data *models.LtGift) error{
	_,err := d.engine.Insert(data)
	return err
}