package orm

import (
	"errors"
	"reflect"
	"time"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrMissPK missing pk error
	ErrMissPK = errors.New("missed pk value")
)

// mysql dbBaser implementation.
type dbBaseMongo struct {
	dbBase
}

var _ dbBaser = new(dbBaseMongo)

// create new mysql dbBaser.
func newdbBaseMongo() dbBaser {
	b := new(dbBaseMongo)
	b.ins = b
	return b
}

// read related records.
func (d *dbBaseMongo) Find(qs *querySet, mi *modelInfo, cond *Condition, container interface{}, tz *time.Location, cols []string) (i int64, err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

	beego.Info(qs.orders)

	opt := options.Find()
	if len(cols) > 0 {
		projection := bson.M{}
		for _, col := range cols {
			projection[col] = 1
		}
		opt.SetProjection(projection)
	}

	if len(qs.orders) > 0 {
		opt.SetSort(getSort(qs.orders))
	}

	if qs.limit != 0 {
		opt.SetLimit(qs.limit)
	}

	if qs.offset != 0 {
		opt.SetSkip(qs.offset)
	}

	filter := convertCondition(cond)

	if qs != nil && qs.forContext {
		// Do something with content
		if err != nil {
			return 0, err
		}
	} else {
		// Do something without content
		cur, err := col.Find(todo, filter, opt)
		defer cur.Close(todo)
		if err != nil {
			return 0, err
		}
		i, err = convertCur(cur, container)
	}

	return
}

// read related records.
func (d *dbBaseMongo) FindOne(qs *querySet, mi *modelInfo, cond *Condition, container interface{}, tz *time.Location, cols []string) (err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

	opt := options.FindOne()
	if len(cols) > 0 {
		projection := bson.M{}
		for _, col := range cols {
			projection[col] = 1
		}
		opt.SetProjection(projection)
	}

	filter := convertCondition(cond)

	if qs != nil && qs.forContext {
		// Do something with content
		if err != nil {
			return err
		}
	} else {
		// Do something without content
		data, err := col.FindOne(todo, filter, opt).DecodeBytes()
		if err != nil {
			return err
		}
		err = bson.Unmarshal(data, container)
	}

	return
}

func (d *dbBaseMongo) Read(q dbQuerier, mi *modelInfo, ind reflect.Value, container interface{}, tz *time.Location, cols []string) (err error) {
	db := q.(*DB).MDB
	col := db.Collection(mi.table)

	opt := options.FindOne()

	var whereCols []string
	var args []interface{}
	if len(cols) > 0 {
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

	filter := bson.M{}
	for i, p := range whereCols {
		filter[p] = args[i]
	}

	// Do something without content
	data, err := col.FindOne(todo, filter, opt).DecodeBytes()
	if err != nil {
		return err
	}
	err = bson.Unmarshal(data, container)

	return
}

func convertCondition(cond *Condition) (filter bson.M) {
	filter = bson.M{}
	for i, p := range cond.params {

		if p.isCond {
			f := convertCondition(p.cond)
			beego.Info(f)
		} else {
			exprs := p.exprs

			num := len(exprs) - 1
			operator := ""
			if operators[exprs[num]] {
				operator = exprs[num]
				exprs = exprs[:num]
			}

			if operator == "" {
				operator = "eq"
			}

			k, v := getCond(exprs, p.args, operator)
			filter[k] = v

			if i > 0 {
				if p.isOr {
					filter = bson.M{
						"$or": bson.A{filter, bson.M{k: v}},
					}
				} else {
					// where += "AND "
					filter[k] = v
				}
			}

		}

	}
	return
}

func getCond(params []string, args []interface{}, operator string) (k string, v interface{}) {
	r := getCondBson(params, args, operator)
	if len(params) > 0 {
		return params[0], r[params[0]]
	}
	return
}

func getCondBson(params []string, args []interface{}, operator string) (r bson.M) {
	r = bson.M{}
	if len(params) == 0 || len(args) == 0 {
		return
	}

	if len(params) == 1 {
		if len(args) == 1 {
			r[params[0]] = bson.M{"$" + operator: args[0]}
			return
		}
		r[params[0]] = bson.M{"$" + operator: args}
		return
	}
	r[params[0]] = bson.M{"$" + operator: getCondBson(params[1:], args, operator)}
	return
}

func getSort(orders []string) (r bson.M) {
	r = bson.M{}
	if len(orders) == 0 {
		return
	}
	for _, order := range orders {
		if order[0] == '-' {
			order = order[1:]
			r[order] = -1
		} else {
			r[order] = 1
		}
	}

	return
}

func convertCur(cur *mongo.Cursor, v interface{}) (i int64, err error) {
	resultv := reflect.ValueOf(v)
	if resultv.Kind() != reflect.Ptr {
		panic("result argument must be a slice address")
	}
	slicev := resultv.Elem()

	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	if slicev.Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	for cur.Next(todo) {
		elemp := reflect.New(elemt)
		if err = bson.Unmarshal(cur.Current, elemp.Interface()); err != nil {
			return
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		i++
	}
	resultv.Elem().Set(slicev.Slice(0, int(i)))
	return
}
