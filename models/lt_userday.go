package models

type LtUserday struct {
	Id         int `xorm:"not null pk autoincr comment('主键') INT(11)"`
	Uid        int `xorm:"INT(11)"`
	Day        int `xorm:"INT(11)"`
	Num        int `xorm:"INT(11)"`
	SysCreated int `xorm:"INT(11)"`
	SysUpdated int `xorm:"INT(11)"`
}
