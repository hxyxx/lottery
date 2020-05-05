package dao

import (
	"github.com/go-xorm/xorm"
	"hxyxx/lottery/models"
	"log"
)

type CodeDao struct {
	engine *xorm.Engine
}

func NewCodeDao(engine *xorm.Engine) *CodeDao{
	return &CodeDao{engine: engine}
}

func (d *CodeDao) Get(id int) *models.LtCode{
	//构建一个model对象
	data := &models.LtCode{Id :id}
	//查询
	ok,err := d.engine.Get(data)
	if ok && err !=nil{
		return data
	}else {
		data.Id = 0
		return data
	}
}
func (d *CodeDao) GetAll() []models.LtCode{
	datalist := make([]models.LtCode,0)
	err := d.engine.
		Desc("id").
		Find(datalist)
	if err != nil {
		log.Println("gift_dao.GetAll error = ",err)
	}
	return datalist
}
func (d *CodeDao) CountAll() int64 {
	num,err := d.engine.Count(&models.LtCode{})
	if err !=nil{
		log.Println("gift_dao.CountAll error : ",err)
		return 0
	}
	return num
}
func (d *CodeDao) Delete(id int) error{
	data := &models.LtCode{Id:id}
	_,err := d.engine.Delete(data)
	return err
}
func (d *CodeDao) update(data *models.LtCode,columns []string) error{
	//MustCols方法指定对象属性即使为空也进行更新，指定字段
	_,err := d.engine.ID(data.Id).MustCols(columns...).Update(data)
	return err
}
func (d *CodeDao) create(data *models.LtCode) error{
	_,err := d.engine.Insert(data)
	return err
}