package orm

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/astaxie/beego"
)

var (
	// ErrMissPK missing pk error
	ErrMissPK = errors.New("missed pk value")
)

// mysql dbBaser implementation.
type dbBaseMysql struct {
	dbBase
}

var _ dbBaser = new(dbBaseMysql)

// create new mysql dbBaser.
func newdbBaseMongo() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}

// read related records.
func (d *dbBaseMysql) ReadBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, container interface{}, tz *time.Location, cols []string) (i int64, err error) {
	val := reflect.ValueOf(container)
	ind := reflect.Indirect(val)

	errTyp := true
	one := true
	isPtr := true

	if val.Kind() == reflect.Ptr {
		fn := ""
		if ind.Kind() == reflect.Slice {
			one = false
			typ := ind.Type().Elem()
			switch typ.Kind() {
			case reflect.Ptr:
				fn = getFullName(typ.Elem())
			case reflect.Struct:
				isPtr = false
				fn = getFullName(typ)
			}
		} else {
			fn = getFullName(ind.Type())
		}
		errTyp = fn != mi.fullName
	}

	if errTyp {
		if one {
			panic(fmt.Errorf("wrong object type `%s` for rows scan, need *%s", val.Type(), mi.fullName))
		} else {
			panic(fmt.Errorf("wrong object type `%s` for rows scan, need *[]*%s or *[]%s", val.Type(), mi.fullName, mi.fullName))
		}
	}

	rlimit := qs.limit
	offset := qs.offset

	var tCols []string
	if len(cols) > 0 {
		hasRel := len(qs.related) > 0 || qs.relDepth > 0
		tCols = make([]string, 0, len(cols))
		var maps map[string]bool
		if hasRel {
			maps = make(map[string]bool)
		}
		for _, col := range cols {
			if fi, ok := mi.fields.GetByAny(col); ok {
				tCols = append(tCols, fi.column)
				if hasRel {
					maps[fi.column] = true
				}
			} else {
				return 0, fmt.Errorf("wrong field/column name `%s`", col)
			}
		}
		if hasRel {
			for _, fi := range mi.fields.fieldsDB {
				if fi.fieldType&IsRelField > 0 {
					if !maps[fi.column] {
						tCols = append(tCols, fi.column)
					}
				}
			}
		}
	} else {
		tCols = mi.fields.dbcols
	}

	colsNum := len(tCols)
	beego.Info(rlimit, offset, colsNum)
	beego.Info(tCols)
	// db := q.GetDB()
	// col := db.Collection(mi.name)

	if qs != nil && qs.forContext {
		// Do something with content
		if err != nil {
			return 0, err
		}
	} else {
		// Do something without content
		// rs, err = q.Query(query, args...)
		if err != nil {
			return 0, err
		}
	}

	refs := make([]interface{}, colsNum)
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	slice := ind
	var cnt int64
	if isPtr {
		// slice = reflect.Append(slice, mind.Addr())
	} else {
		// slice = reflect.Append(slice, mind)
	}

	if !one {
		if cnt > 0 {
			ind.Set(slice)
		} else {
			// when a result is empty and container is nil
			// to set a empty container
			if ind.IsNil() {
				ind.Set(reflect.MakeSlice(ind.Type(), 0, 0))
			}
		}
	}
	return
}

func (d *dbBaseMysql) Read(q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location, cols []string) (err error) {
	var whereCols []string
	var args []interface{}
	if len(cols) > 0 {
		var err error
		whereCols = make([]string, 0, len(cols))
		args, _, err = d.collectValues(mi, ind, cols, false, false, &whereCols, tz)
		if err != nil {
			return err
		}
	} else {
		// default use pk value as where condtion.
		pkColumn, pkValue, ok := getExistPk(mi, ind)
		if !ok {
			return ErrMissPK
		}
		whereCols = []string{pkColumn}
		args = append(args, pkValue)
	}

	colsNum := len(mi.fields.dbcols)

	beego.Info(colsNum)

	refs := make([]interface{}, colsNum)
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}
	// 把查询结果写入refs
	mdb := q.GetDB()
	beego.Info(mdb)

	elm := reflect.New(mi.addrField.Elem().Type())
	mind := reflect.Indirect(elm)
	d.setColsValues(mi, &mind, mi.fields.dbcols, refs, tz)
	ind.Set(mind)
	return
}
