package orm

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoIndexKey = errors.New("have not a index key")
)

type Index struct {
	Keys []string
	options.IndexOptions
}
type indexView struct {
	index mongo.IndexView
}

var _ IndexViewer = new(indexView)

// list all index
func (iv *indexView) List() (val interface{}, err error) {
	res := []map[string]interface{}{}
	opt := options.ListIndexes()
	cur, err := iv.index.List(todo, opt)
	if err != nil {
		return
	}
	defer cur.Close(todo)
	err = cur.All(todo, &res)
	return res, err
}

// create one index by indexModel
func (iv *indexView) CreateOne(index Index, t ...time.Duration) (id string, err error) {
	opts := options.CreateIndexes()
	if len(t) > 0 {
		opts.SetMaxTime(t[0] * time.Second)
	}

	keys, iopts, err := convertIndex(index)

	if err != nil {
		return
	}

	model := mongo.IndexModel{
		Keys:    keys,
		Options: iopts,
	}

	return iv.index.CreateOne(todo, model, opts)
}

// creat many index by indexModels
func (iv *indexView) CreateMany(indexs []Index, t ...time.Duration) (ids []string, err error) {
	opts := options.CreateIndexes()
	if len(t) > 0 {
		opts.SetMaxTime(t[0] * time.Second)
	}
	models := []mongo.IndexModel{}
	for _, index := range indexs {
		keys, iopts, err1 := convertIndex(index)
		if err1 != nil {
			err = err1
			return
		}
		model := mongo.IndexModel{
			Keys:    keys,
			Options: iopts,
		}
		models = append(models, model)
	}

	return iv.index.CreateMany(todo, models, opts)
}

// drop one index by index name
func (iv *indexView) DropOne(name string, t ...time.Duration) (err error) {
	opts := options.DropIndexes()
	if len(t) > 0 {
		opts.SetMaxTime(t[0] * time.Second)
	}

	_, err = iv.index.DropOne(todo, name, opts)
	return
}

// drop all index
func (iv *indexView) DropAll(t ...time.Duration) (err error) {
	opts := options.DropIndexes()
	if len(t) > 0 {
		opts.SetMaxTime(t[0] * time.Second)
	}

	_, err = iv.index.DropAll(todo, opts)
	return
}

// new indexView
func newIndexView(iv mongo.IndexView) IndexViewer {
	v := new(indexView)
	v.index = iv
	return v
}

func convertIndex(index Index) (keys bson.M, iopts *options.IndexOptions, err error) {
	if len(index.Keys) < 1 {
		err = ErrNoIndexKey
		return
	}

	keys = bson.M{}

	for _, v := range index.Keys {
		if v[0] == '-' {
			v = v[1:]
			keys[v] = -1
		} else {
			keys[v] = 1
		}
	}

	iopts = &index.IndexOptions
	return
}
