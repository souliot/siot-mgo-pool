package orm

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrMissPK missing pk error
	ErrMissPK = errors.New("missed pk value")
)

type OperatorUpdate string

var (
	MgoSet         OperatorUpdate = "$set"
	MgoUnSet       OperatorUpdate = "$unset"
	MgoInc         OperatorUpdate = "$inc"
	MgoPush        OperatorUpdate = "$push"
	MgoPushAll     OperatorUpdate = "$pushAll"
	MgoAddToSet    OperatorUpdate = "$addToSet"
	MgoPop         OperatorUpdate = "$pop"
	MgoPull        OperatorUpdate = "$pull"
	MgoPullAll     OperatorUpdate = "$pullAll"
	MgoRename      OperatorUpdate = "$rename"
	MgoSetOnInsert OperatorUpdate = "$setOnInsert"
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

// read one record.
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

	if len(qs.orders) > 0 {
		opt.SetSort(getSort(qs.orders))
	}

	if qs.offset != 0 {
		opt.SetSkip(qs.offset)
	}

	filter := convertCondition(cond)

	if qs != nil && qs.forContext {
		err = col.FindOne(qs.ctx, filter, opt).Decode(container)
	} else {
		err = col.FindOne(todo, filter, opt).Decode(container)
	}

	return
}

// read one record.
func (d *dbBaseMongo) Distinct(qs *querySet, mi *modelInfo, cond *Condition, tz *time.Location, field string) (res []interface{}, err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)
	opt := options.Distinct()

	filter := convertCondition(cond)

	if qs != nil && qs.forContext {
		return col.Distinct(qs.ctx, field, filter, opt)
	} else {
		return col.Distinct(todo, field, filter, opt)
	}
}

// read all records.
func (d *dbBaseMongo) Find(qs *querySet, mi *modelInfo, cond *Condition, container interface{}, tz *time.Location, cols []string) (err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

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
	cur := &mongo.Cursor{}
	if qs != nil && qs.forContext {
		// Do something with content
		cur, err = col.Find(qs.ctx, filter, opt)
	} else {
		// Do something without content
		cur, err = col.Find(todo, filter, opt)
	}
	// defer cur.Close(todo)
	if err != nil {
		return
	}
	err = cur.All(todo, container)

	return
}

// get the recodes count.
func (d *dbBaseMongo) Count(qs *querySet, mi *modelInfo, cond *Condition, tz *time.Location) (i int64, err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

	opt := options.Count()

	filter := convertCondition(cond)

	if qs != nil && qs.forContext {
		if len(filter) == 0 {
			i, err = col.EstimatedDocumentCount(qs.ctx, nil)
		} else {
			i, err = col.CountDocuments(qs.ctx, filter, opt)
		}
	} else {
		// Do something without content
		if len(filter) == 0 {
			i, err = col.EstimatedDocumentCount(todo, nil)
		} else {
			i, err = col.CountDocuments(todo, filter, opt)
		}
	}

	return
}

// update the recodes.
func (d *dbBaseMongo) UpdateMany(qs *querySet, mi *modelInfo, cond *Condition, operator OperatorUpdate, params Params, tz *time.Location) (i int64, err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

	opt := options.Update()

	filter := convertCondition(cond)
	update := bson.M{}
	for col, val := range params {
		// if fi, ok := mi.fields.GetByAny(col); !ok || !fi.dbcol {
		// 	panic(fmt.Errorf("wrong field/column name `%s`", col))
		// } else {
		update[col] = val
		// }
	}
	update = bson.M{
		string(operator): update,
	}
	r := &mongo.UpdateResult{}
	if qs != nil && qs.forContext {
		r, err = col.UpdateMany(qs.ctx, filter, update, opt)
	} else {
		// Do something without content
		r, err = col.UpdateMany(todo, filter, update, opt)
	}
	if err != nil {
		return
	}

	i = r.ModifiedCount

	return
}

// delete the recodes.
func (d *dbBaseMongo) DeleteMany(qs *querySet, mi *modelInfo, cond *Condition, tz *time.Location) (i int64, err error) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

	opt := options.Delete()

	filter := convertCondition(cond)

	r := &mongo.DeleteResult{}
	if qs != nil && qs.forContext {
		r, err = col.DeleteMany(qs.ctx, filter, opt)
	} else {
		// Do something without content
		r, err = col.DeleteMany(todo, filter, opt)
	}
	if err != nil {
		return
	}

	i = r.DeletedCount

	return
}

// get indexview.
func (d *dbBaseMongo) Indexes(qs *querySet, mi *modelInfo, tz *time.Location) (iv IndexViewer) {
	db := qs.orm.db.(*DB).MDB
	col := db.Collection(mi.table)

	return newIndexView(col.Indexes())
}

// read one record.
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

// insert one record.
func (d *dbBaseMongo) InsertOne(q dbQuerier, mi *modelInfo, ind reflect.Value, container interface{}, tz *time.Location) (id interface{}, err error) {
	db := q.(*DB).MDB
	col := db.Collection(mi.table)
	_, _, b := getExistPk(mi, ind)
	name := mi.fields.pk.name

	if !b {
		reflect.ValueOf(container).Elem().FieldByName(name).SetString(primitive.NewObjectID().Hex())
	}

	opt := options.InsertOne()

	// Do something without content
	data, err := col.InsertOne(todo, container, opt)
	if err != nil {
		return
	}
	id = data.InsertedID
	return
}

// insert all records.
func (d *dbBaseMongo) InsertMany(q dbQuerier, mi *modelInfo, ind reflect.Value, containers interface{}, tz *time.Location) (ids interface{}, err error) {
	db := q.(*DB).MDB
	col := db.Collection(mi.table)
	_, _, b := getExistPk(mi, ind)
	name := mi.fields.pk.name
	sind := reflect.Indirect(reflect.ValueOf(containers))

	cs := []interface{}{}

	if !b {
		for i := 0; i < sind.Len(); i++ {
			c := reflect.Indirect(sind.Index(i))
			c.FieldByName(name).SetString(primitive.NewObjectID().Hex())
			cs = append(cs, c.Interface())
		}
	}

	opt := options.InsertMany()

	// Do something without content
	data, err := col.InsertMany(todo, cs, opt)
	if err != nil {
		return
	}
	ids = data.InsertedIDs
	return
}

// update one record.
func (d *dbBaseMongo) UpdateOne(q dbQuerier, mi *modelInfo, ind reflect.Value, container interface{}, tz *time.Location, cols []string) (id interface{}, err error) {
	db := q.(*DB).MDB
	col := db.Collection(mi.table)
	c, val, b := getExistPk(mi, ind)
	if !b {
		return nil, ErrHaveNoPK
	}

	opt := options.Update()
	var whereCols []string
	var args []interface{}

	if len(cols) == 0 {
		cols = mi.fields.dbcols
	}
	whereCols = make([]string, 0, len(cols))
	args, _, err = d.collectValues(mi, ind, cols, false, false, &whereCols, tz)
	if err != nil {
		return
	}

	filter := bson.M{
		c: val,
	}

	update := bson.M{}
	for i, p := range whereCols {
		if p != c {
			update[p] = args[i]
		}
	}

	update = bson.M{
		"$set": update,
	}

	// Do something without content
	data, err := col.UpdateOne(todo, filter, update, opt)
	id = data.UpsertedID
	return
}

// delete one record.
func (d *dbBaseMongo) DeleteOne(q dbQuerier, mi *modelInfo, ind reflect.Value, container interface{}, tz *time.Location, cols []string) (cnt interface{}, err error) {
	db := q.(*DB).MDB
	col := db.Collection(mi.table)

	opt := options.Delete()
	var whereCols []string
	var args []interface{}
	if len(cols) > 0 {
		whereCols = make([]string, 0, len(cols))
		args, _, err = d.collectValues(mi, ind, cols, false, false, &whereCols, tz)
		if err != nil {
			return
		}
	} else {
		// default use pk value as where condtion.
		pkColumn, pkValue, ok := getExistPk(mi, ind)
		if !ok {
			return nil, ErrMissPK
		}
		whereCols = []string{pkColumn}
		args = append(args, pkValue)
	}

	filter := bson.M{}
	for i, p := range whereCols {
		filter[p] = args[i]
	}

	// Do something without content
	data, err := col.DeleteOne(todo, filter, opt)
	cnt = data.DeletedCount
	return
}

func convertCondition(cond *Condition) (filter bson.M) {
	filter = bson.M{}
	if cond == nil {
		return
	}
	for i, p := range cond.params {
		if p.isCond {
			f := convertCondition(p.cond)
			if i > 0 {
				if p.isOr {
					filter = bson.M{
						"$or": bson.A{filter, f},
					}
				} else {
					// where += "AND "
					filter = bson.M{
						"$and": bson.A{filter, f},
					}
				}
			} else {
				filter = f
			}

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

			if i > 0 {
				if p.isOr {
					filter = bson.M{
						"$or": bson.A{filter, bson.M{k: v}},
					}
				} else {
					// where += "AND "
					filter[k] = v
				}
			} else {
				filter[k] = v
			}

		}

	}
	return
}

func getCond(params []string, args []interface{}, operator string) (k string, v interface{}) {
	k = strings.Join(params, ".")
	if len(args) == 0 {
		v = bson.M{}
	} else if len(args) == 1 {
		v = bson.M{
			"$" + operator: args[0],
		}
	} else {
		v = bson.M{
			"$" + operator: args,
		}
	}

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
