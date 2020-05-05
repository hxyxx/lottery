package models

type LtResult struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	GiftId     int    `xorm:"comment('奖品id，关联prize表') INT(11)"`
	GiftName   int    `xorm:"comment('奖品名称') INT(11)"`
	GiftType   int    `xorm:"comment('奖品类型，同coupon.type') INT(11)"`
	Uid        int    `xorm:"comment('用户id') INT(11)"`
	Username   string `xorm:"comment('用户名') VARCHAR(50)"`
	PrizeCode  int    `xorm:"comment('抽奖编码(4位随机数)') INT(11)"`
	GiftData   string `xorm:"comment('获奖信息') VARCHAR(255)"`
	SysCreated int    `xorm:"comment('创建时间') INT(11)"`
	SysIp      string `xorm:"comment('用户抽奖的ip') VARCHAR(50)"`
	SysStatus  int    `xorm:"comment('状态，0正常，1删除，2作弊') SMALLINT(6)"`
}