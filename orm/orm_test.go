package orm

import (
	"testing"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
)

type Logs struct {
	Id       int    `bson:"_id"`
	Ltype    string `orm:"column(type)" bson:"type"`
	UserName string `orm:"column(username)"  bson:"username"`
	L2       *Logs2 `orm:"rel(one)" bson:"l2"`
}

type Logs2 struct {
	Id       int    `bson:"_id"`
	Ltype    string `orm:"column(type)" bson:"type"`
	UserName string `orm:"column(username)" bson:"username"`
}

func (m *Logs) TableName() string {
	return "log"
}

func init() {
	beego.SetLogFuncCall(true)
	RegisterModel(new(Logs))
	RegisterModel(new(Logs2))
}

func TestDB(t *testing.T) {

	RegisterDriver("mongo", DRMongo)
	RegisterDataBase("default", "mongo", "mongodb://yapi:abcd1234@192.168.50.200:27017/yapi")

	o := NewOrm()
	o.Using("default")
	var ls []Logs
	qs := o.QueryTable("log")
	num, err := qs.Filter("username", "admin").Filter("_id__in", 16, 11, 12).OrderBy("-_id", "type").Offset(0).Limit(100).All(&ls, "username", "type")
	beego.Info(num, err)
	beego.Info(ls)

	l := Logs{
		UserName: "admin",
		Ltype:    "group",
	}
	err = o.Read(&l, "UserName", "Ltype")
	beego.Info(l, err == mongo.ErrNoDocuments)
}
