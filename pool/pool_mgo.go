package pool

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func RegisterMgoPool(poolName string, url string, force bool, params ...int) (err error) {
	//factory 创建连接的方法
	factory := func() (interface{}, error) {
		return mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	}

	//close 关闭连接的方法
	close := func(v interface{}) error {
		return v.(*mongo.Client).Disconnect(context.Background())
	}

	//ping 检测连接的方法
	ping := func(v interface{}) error {
		return v.(*mongo.Client).Ping(context.Background(), readpref.Primary())
	}

	var (
		size int           = 5
		cap  int           = 20
		idle time.Duration = 30
	)

	for i, v := range params {
		switch i {
		case 0:
			size = v
		case 1:
			cap = v
		case 2:
			idle = time.Duration(v)
		}
	}

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &Config{
		InitialCap: size,
		MaxCap:     cap,
		Factory:    factory,
		Close:      close,
		Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: idle * time.Second,
	}
	mgoPool, err := NewChannelPool(poolConfig)

	if !pools.add(poolName, mgoPool, force) {
		return ErrRegisterPool
	}
	return
}

func GetMgoClient(poolName string) (c *mongo.Client, err error) {
	if p, ok := pools.get(poolName); ok {
		v, err := p.Get()
		if err == nil {
			c = v.(*mongo.Client)
			defer PutMgoClient(poolName, c)
		}
		return c, err
	}
	return nil, ErrGetConnection
}

func PutMgoClient(poolName string, c *mongo.Client) (err error) {
	if p, ok := pools.get(poolName); ok {
		err = p.Put(c)
		return
	}
	return ErrPutConnection
}
