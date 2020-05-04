package models

type LtGift struct {
	Id           int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	Title        string `xorm:"comment('奖品名称') VARCHAR(255)"`
	PrizeNum     int    `xorm:"comment('奖品数量，0 无限量，>0限量， <0无奖品') INT(11)"`
	LeftNum      int    `xorm:"comment('剩余奖品数量') INT(11)"`
	GiftCode     string `xorm:"comment('0-9999表示100%,0-0表示1%') index VARCHAR(50)"`
	GiftTime     int    `xorm:"comment('发奖周期D天') INT(11)"`
	Img          string `xorm:"comment('奖品图片') VARCHAR(255)"`
	Displayorder int    `xorm:"comment('位置序号，小的排在前边') INT(11)"`
	Gtype        int    `xorm:"comment('奖品类型，0 虚拟币，1虚拟券，2实物-小奖，3实物-大奖') INT(11)"`
	Gdata        string `xorm:"comment('扩展数据，如虚拟币数量') VARCHAR(255)"`
	TimeBegin    int    `xorm:"comment('开始时间') INT(11)"`
	TimeEnd      int    `xorm:"comment('结束时间') INT(11)"`
	PrizeData    string `xorm:"comment('发奖计划[时间1，数量1],[时间2，数量2]') MEDIUMTEXT"`
	PrizeBegin   int    `xorm:"comment('发奖周期的开始') INT(11)"`
	PrizeEnd     int    `xorm:"comment('发奖周期的结束') INT(11)"`
	SysStatus    int    `xorm:"comment('状态，0正常，1删除') SMALLINT(6)"`
	SysCreated   int    `xorm:"comment('创建时间') INT(11)"`
	SysUpdated   int    `xorm:"comment('修改时间') INT(11)"`
	SysIp        string `xorm:"comment('操作人ip') VARCHAR(50)"`
}
