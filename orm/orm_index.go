package orm

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func (iv *indexView) CreateOne(model mongo.IndexModel, t ...time.Duration) (id string, err error) {
	opts := options.CreateIndexes()
	if len(t) > 0 {
		opts.SetMaxTime(t[0] * time.Second)
	}

	return iv.index.CreateOne(todo, model, opts)
}

// creat many index by indexModels
func (iv *indexView) CreateMany(models []mongo.IndexModel, t ...time.Duration) (ids []string, err error) {
	opts := options.CreateIndexes()
	if len(t) > 0 {
		opts.SetMaxTime(t[0] * time.Second)
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
