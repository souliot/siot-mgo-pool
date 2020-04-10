package orm

import (
	"testing"
	"time"

	"github.com/astaxie/beego/logs"
)

type Logs struct {
	Ids      string `bson:"_id"`
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
	logs.SetLogFuncCall(true)

	RegisterModel(new(Logs))
	RegisterModel(new(Logs2))

	RegisterDriver("mongo", DRMongo)
	RegisterDataBase("default", "mongo", "mongodb://yapi:abcd1234@192.168.50.200:27017/yapi")
	// RegisterDataBase("default", "mongo", "mongodb://yapi:abcd1234@vm:27017/yapi")

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
	logs.Info(err, l)
}
func TestReadOrCreate(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	c, id, err := o.ReadOrCreate(&l, "UserName")
	logs.Info(c, id, err, l)
}

func TestInsert(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	id, err := o.Insert(&l)
	logs.Info(id, err)
}

func TestInsertMulti(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	ls := []Logs{}
	ls = append(ls, l)
	ls = append(ls, l)
	id, err := o.InsertMulti(ls)
	logs.Info(id, err)
}

func TestUpdate(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	l.Ids = "5e7431f78c1b4111312cce2d"
	l.Ltype = "group3"
	id, err := o.Update(&l, "Ltype")
	logs.Info(id, err, l)
}

func TestDelete(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	l.Ids = "5e71816fee8b0d2ba0d24939"
	l.Ltype = "group3"
	cnt, err := o.Delete(&l, "Ltype")
	logs.Info(cnt, err)
}

func TestQsOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	err := qs.Filter("username", "linleizhou1234").One(&l, "username", "type")
	logs.Info(err, l)
}

func TestQsAll(t *testing.T) {
	o := NewOrm()
	o.Using("default")
	var ls []Logs
	qs := o.QueryTable("log")
	err := qs.Filter("username__regex", "linleizhou").OrderBy("-_id", "Ltype").Limit(2, 0).All(&ls)
	// num, err := qs.All(&ls)
	logs.Info(err)
	logs.Info(ls)
}

func TestQsCount(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	num, err := qs.Filter("username", "linleizhou1234").Count()
	// num, err := qs.Count()
	logs.Info(num, err)
}
func TestQsUpdate(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	num, err := qs.Filter("_id", "5e7431f78c1b4111312cce2d").Update(MgoSet, Params{
		"type": "group",
	})
	logs.Info(num, err)
}
func TestQsDelete(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	num, err := qs.Filter("type", "group3").Delete()
	logs.Info(num, err)
}
func TestQsIndexList(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	indexes, err := qs.IndexView().List()
	logs.Info(indexes, err)
}
func TestQsIndexCreateOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")
	qs := o.QueryTable("log")

	index := Index{}
	index.Keys = []string{"-username", "_id"}
	index.SetName("username").SetUnique(true)

	indexes, err := qs.IndexView().CreateOne(index)
	logs.Info(indexes, err)

}
func TestQsIndexDropOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")

	qs := o.QueryTable("log")
	err := qs.IndexView().DropOne("username")
	logs.Info(err)
}
func TestOther(t *testing.T) {
	// uri := "mongodb://@192.168.0.4:27017/Darwin-XYY"
	// cs, err := connstring.Parse(uri)
	// logs.Info(err)
	// logs.Info(cs.Database)
	logs.Info(time.Now().Unix())
	logs.Info(time.Now().UTC().Unix())
}
