package orm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type orm struct {
	alias *alias
	isTx  bool
	db    dbQuerier
}

// 下划线用来判断结构体是否实现了接口，
// 如果没有实现，在编译的时候就能暴露出问题，
// 如果没有这个判断，后代码中使用结构体没有实现的接口方法，在编译器是不会报错的。

var (
	Debug            = false
	DebugLog         = NewLog(os.Stdout)
	DefaultRowsLimit = -1
	DefaultRelsDepth = 2
	DefaultTimeLoc   = time.Local
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> err no rows")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")
	ErrHaveNoPK      = errors.New("<Ormer> the PK value should not be nil")
	ErrNoDocuments   = mongo.ErrNoDocuments
	todo             = context.TODO()
)

// Params stores the Params
type Params map[string]interface{}

// ParamsList stores paramslist
type ParamsList []interface{}

const (
	formatTime     = "15:04:05"
	formatDate     = "2006-01-02"
	formatDateTime = "2006-01-02 15:04:05"
)

var _ Ormer = new(orm)

// get model info and model reflect value
func (o *orm) getMiInd(md interface{}, needPtr bool) (mi *modelInfo, ind reflect.Value) {
	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ := ind.Type()
	if needPtr && val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Ormer> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}
	name := getFullName(typ)
	if mi, ok := modelCache.getByFullName(name); ok {
		return mi, ind
	}
	panic(fmt.Errorf("<Ormer> table: `%s` not found, make sure it was registered with `RegisterModel()`", name))
}

// read data to model
func (o *orm) Read(md interface{}, cols ...string) (err error) {
	mi, ind := o.getMiInd(md, true)
	return o.alias.DbBaser.Read(o.db, mi, ind, md, o.alias.TZ, cols)
}

// Try to read a row from the database, or insert one if it doesn't exist
func (o *orm) ReadOrCreate(md interface{}, col1 string, cols ...string) (created bool, id interface{}, err error) {
	cols = append([]string{col1}, cols...)
	mi, ind := o.getMiInd(md, true)
	err = o.alias.DbBaser.Read(o.db, mi, ind, md, o.alias.TZ, cols)
	if err == mongo.ErrNoDocuments {
		// Create
		id, err = o.Insert(md)
		return (err == nil), id, err
	}

	vid := ind.FieldByIndex(mi.fields.pk.fieldIndex)
	if mi.fields.pk.fieldType&IsPositiveIntegerField > 0 {
		id = int64(vid.Uint())
	} else if mi.fields.pk.rel {
		return o.ReadOrCreate(vid.Interface(), mi.fields.pk.relModelInfo.fields.pk.name)
	}

	return false, id, err
}

// insert model data to database
func (o *orm) Insert(md interface{}) (id interface{}, err error) {
	mi, ind := o.getMiInd(md, true)
	id, err = o.alias.DbBaser.InsertOne(o.db, mi, ind, md, o.alias.TZ)
	return
}

// insert models data to database
func (o *orm) InsertMulti(mds interface{}) (ids interface{}, err error) {
	sind := reflect.Indirect(reflect.ValueOf(mds))
	switch sind.Kind() {
	case reflect.Array, reflect.Slice:
		if sind.Len() == 0 {
			return nil, ErrArgs
		}
	default:
		return nil, ErrArgs
	}
	ind := reflect.Indirect(sind.Index(0))
	mi, _ := o.getMiInd(ind.Interface(), false)
	ids, err = o.alias.DbBaser.InsertMany(o.db, mi, ind, mds, o.alias.TZ)
	return
}

// cols set the columns those want to update.
func (o *orm) Update(md interface{}, cols ...string) (interface{}, error) {
	mi, ind := o.getMiInd(md, true)
	return o.alias.DbBaser.UpdateOne(o.db, mi, ind, md, o.alias.TZ, cols)
}

// delete model in database
// cols shows the delete conditions values read from. default is pk
func (o *orm) Delete(md interface{}, cols ...string) (interface{}, error) {
	mi, ind := o.getMiInd(md, true)
	return o.alias.DbBaser.DeleteOne(o.db, mi, ind, md, o.alias.TZ, cols)
}

// set auto pk field
func (o *orm) setPk(mi *modelInfo, ind reflect.Value, id int64) {
	if mi.fields.pk.auto {
		if mi.fields.pk.fieldType&IsPositiveIntegerField > 0 {
			ind.FieldByIndex(mi.fields.pk.fieldIndex).SetUint(uint64(id))
		} else {
			ind.FieldByIndex(mi.fields.pk.fieldIndex).SetInt(id)
		}
	}
}

// return a QuerySeter for table operations.
// table name can be string or struct.
// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
func (o *orm) QueryTable(ptrStructOrTableName interface{}) (qs QuerySeter) {
	var name string
	if table, ok := ptrStructOrTableName.(string); ok {
		name = nameStrategyMap[MongoNameStrategy](table)
		if mi, ok := modelCache.get(name); ok {
			qs = newQuerySet(o, mi)
		}
	} else {
		name = getFullName(indirectType(reflect.TypeOf(ptrStructOrTableName)))
		if mi, ok := modelCache.getByFullName(name); ok {
			qs = newQuerySet(o, mi)
		}
	}
	if qs == nil {
		panic(fmt.Errorf("<Ormer.QueryTable> table name: `%s` not exists", name))
	}
	return
}

func NewOrm() Ormer {
	BootStrap() // execute only once

	o := new(orm)
	err := o.Using("default")
	if err != nil {
		panic(err)
	}
	return o
}

func (o *orm) Using(name string) error {
	if o.isTx {
		panic(fmt.Errorf("<Ormer.Using> transaction has been start, cannot change db"))
	}
	if al, ok := dataBaseCache.get(name); ok {
		o.alias = al
		db, err := al.getDB()
		if err != nil {
			return err
		}
		o.db = db
	} else {
		return fmt.Errorf("<Ormer.Using> unknown db alias name `%s`", name)
	}
	return nil
}
func (o *orm) Begin() (err error) {
	if o.isTx {
		return
	}

	err = o.db.Begin()
	if err != nil {
		return err
	}
	o.isTx = true
	return
}

func (o *orm) Commit() (err error) {
	if !o.isTx {
		return ErrTxDone
	}
	err = o.db.Commit()
	if err == nil {
		o.isTx = false
		o.Using(o.alias.Name)
	} else {
		return ErrTxDone
	}
	return
}
func (o *orm) Rollback() (err error) {
	if !o.isTx {
		return ErrTxDone
	}
	err = o.db.Rollback()
	if err == nil {
		o.isTx = false
		o.Using(o.alias.Name)
	} else {
		return ErrTxDone
	}
	return
}
