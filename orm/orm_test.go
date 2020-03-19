package orm

import (
	"testing"

	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/astaxie/beego"
)

type Logs struct {
	Ids      string `orm:"column(_id)" bson:"_id"`
	Ltype    string `orm:"column(type)" bson:"type"`
	UserName string `orm:"column(username)"  bson:"username"`
	L2       *Logs2 `bson:"l2"`
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

	RegisterDriver("mongo", DRMongo)
	// RegisterDataBase("default", "mongo", "mongodb://yapi:abcd1234@192.168.50.200:27017/yapi")
	RegisterDataBase("default", "mongo", "mongodb://yapi:abcd1234@vm:27017/yapi")

}

var (
	l = Logs{
		UserName: "linleizhou1234",
		Ltype:    "group",
	}
)

func TestRead(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	err := o.Read(&l, "UserName")
	beego.Info(err, l)
}
func TestReadOrCreate(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	c, id, err := o.ReadOrCreate(&l, "UserName")
	beego.Info(c, id, err, l)
}

func TestInsert(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	id, err := o.Insert(&l)
	beego.Info(id, err)
}

func TestInsertMulti(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	ls := []Logs{}
	ls = append(ls, l)
	ls = append(ls, l)
	id, err := o.InsertMulti(ls)
	beego.Info(id, err)
}

func TestUpdate(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	l.Ids = "5e72fce41465edf903db7a62"
	l.Ltype = "group3"
	id, err := o.Update(&l, "Ltype")
	beego.Info(id, err, l)
}

func TestDelete(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	l.Ids = "5e71816fee8b0d2ba0d24939"
	l.Ltype = "group3"
	cnt, err := o.Delete(&l, "Ltype")
	beego.Info(cnt, err)
}

func TestQsOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	err := qs.Filter("username", "linleizhou1234").One(&l, "username", "type")
	beego.Info(err, l)
}

func TestQsAll(t *testing.T) {
	o := NewOrm()
	o.Using("default")
	var ls []Logs
	qs := o.QueryTable("log")
	num, err := qs.Filter("username", "linleizhou1234").OrderBy("-_id", "Ltype").Offset(0).Limit(100).All(&ls)
	beego.Info(num, err)
	beego.Info(ls)
}

func TestQsCount(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	num, err := qs.Filter("username", "linleizhou1234").Count()
	beego.Info(num, err)
}
func TestQsUpdate(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	num, err := qs.Filter("_id", "5e72fce41465edf903db7a63").Update(Params{
		"type": "group3",
	})
	beego.Info(num, err)
}
func TestQsDelete(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	num, err := qs.Filter("type", "group3").Delete()
	beego.Info(num, err)
}
func TestQsIndexList(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	indexes, err := qs.IndexView().List()
	beego.Info(indexes, err)
}
func TestQsIndexCreateOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")
	qs := o.QueryTable("log")

	index := Index{}
	index.Keys = []string{"-username", "_id"}
	index.SetName("username").SetUnique(true)

	indexes, err := qs.IndexView().CreateOne(index)
	beego.Info(indexes, err)

}
func TestQsIndexDropOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	err := qs.IndexView().DropOne("username")
	beego.Info(err)
}
func TestOther(t *testing.T) {
	uri := "mongodb://@192.168.0.4:27017/Darwin-XYY"
	cs, err := connstring.Parse(uri)
	beego.Info(err)
	beego.Info(cs.Database)
}
