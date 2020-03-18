package orm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type indexView struct {
	index mongo.IndexView
}

var _ IndexViewer = new(indexView)

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
func (iv *indexView) CreateOne(model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (id string, err error) {
	return iv.index.CreateOne(todo, model, opts...)
}
func (iv *indexView) CreateMany(models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) (ids []string, err error) {
	return iv.index.CreateMany(todo, models, opts...)
}

func newIndexView(iv mongo.IndexView) IndexViewer {
	v := new(indexView)
	v.index = iv
	return v
}
