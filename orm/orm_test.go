package orm

import (
	"reflect"
	"testing"

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
	l.UserName = "1qewe"
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
	l.Ids = "5e70bb453d550128d8493a3a"
	l.Ltype = "group3"
	id, err := o.Update(&l)
	beego.Info(id, err)
}

func TestDelete(t *testing.T) {
	o := NewOrm()
	o.Using("default")
	l.Ids = "5e70bb453d550128d8493a3a"
	l.Ltype = "group"
	id, err := o.Delete(&l, "Ltype")
	beego.Info(id, err)
}

func TestOne(t *testing.T) {
	o := NewOrm()
	o.Using("default")
	qs := o.QueryTable("log")
	err := qs.Filter("username", "linleizhou1234").One(&l, "username", "type")
	beego.Info(err, l)
}

func TestAll(t *testing.T) {
	o := NewOrm()

	o.Using("default")
	var ls []Logs
	qs := o.QueryTable("log")
	num, err := qs.Filter("username", "linleizhou1234").OrderBy("-_id", "type").Offset(0).Limit(100).All(&ls)
	beego.Info(num, err)
	beego.Info(ls)
}

func TestMy(t *testing.T) {
	Up2(&l)
	beego.Info(l)
}

func Update(d interface{}) {
	v := reflect.ValueOf(d).Elem()
	v.FieldByName("Ids").SetString("asdsda")
}

func Up2(d interface{}) {
	Update(d)
}
