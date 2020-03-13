package main

import (
	"github.com/astaxie/beego"
	"github.com/souliot/siot-mgo-pool/orm"
)

type Log struct {
	Id       int    `bson:"_id"`
	Ltype    string `bson:"group"`
	UserName string `orm:"column" bson:"username"`
}

func (m *Log) TableName() string {
	return "Log"
}

func init() {
	beego.SetLogFuncCall(true)
	orm.RegisterModel(new(Log))
}

func main() {
	// pool.RegisterMgoPool("default", "mongodb://yapi:abcd1234@vm:27017/yapi")
	orm.RegisterDriver("mongo", orm.DRMongo)
	orm.RegisterDataBase("default", "mongo", "mongodb://yapi:abcd1234@vm:27017/yapi")

	o := orm.NewOrm()
	o.Using("default")
	l := &Log{16, "group", "linleizhou"}
	qs := o.QueryTable("Log")
	qs.One(l, "UserName")

}
