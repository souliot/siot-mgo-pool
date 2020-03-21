// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"fmt"
	"sync"
	"time"

	"github.com/souliot/siot-mgo-pool/pool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// DriverType database driver constant int.
type DriverType int

// Enum the Database driver
const (
	_ DriverType = iota // int enum type
	DRMongo
)

var (
	dataBaseCache = &_dbCache{cache: make(map[string]*alias)}
	drivers       = map[string]DriverType{
		"mongo": DRMongo,
	}
	dbBasers = map[DriverType]dbBaser{
		DRMongo: newdbBaseMongo(),
	}
)

// database alias cacher.
type _dbCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

// add database alias with original name.
func (ac *_dbCache) add(name string, al *alias) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; !ok {
		ac.cache[name] = al
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dbCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_dbCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}

type DB struct {
	MDB     *mongo.Database
	Session mongo.Session
}

var _ dbQuerier = new(DB)

func (d *DB) Begin() (err error) {
	d.Session, err = d.MDB.Client().StartSession()
	if err != nil {
		return
	}
	defer d.Session.EndSession(todo)

	//开始事务
	err = d.Session.StartTransaction()

	return
}

func (d *DB) Commit() (err error) {
	return d.Session.CommitTransaction(todo)
}
func (d *DB) Rollback() (err error) {
	return d.Session.AbortTransaction(todo)
}

type alias struct {
	Name         string
	Driver       DriverType
	DriverName   string
	DataSource   string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
	DbBaser      dbBaser
	TZ           *time.Location
	Engine       string
}

func (al *alias) getDB() (db *DB, err error) {
	if al.Name == "" {
		al.Name = "default"
	}
	client, err := pool.GetMgoClient(al.Name)
	if err != nil {
		DebugLog.Println(err.Error())
		return
	}

	db = &DB{client.Database(al.DbName), nil}
	return
}

func detectTZ(al *alias) {
	// orm timezone system match database
	// default use Local
	al.TZ = DefaultTimeLoc
}

func addAlias(aliasName, driverName string) (*alias, error) {
	al := new(alias)
	al.Name = aliasName
	al.DriverName = driverName

	if dr, ok := drivers[driverName]; ok {
		al.DbBaser = dbBasers[dr]
		al.Driver = dr
	} else {
		return nil, fmt.Errorf("driver name `%s` have not registered", driverName)
	}

	if !dataBaseCache.add(aliasName, al) {
		return nil, fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
	}

	return al, nil
}

// AddAliasWthDB add a aliasName for the drivename
func AddAlias(aliasName, driverName string) error {
	_, err := addAlias(aliasName, driverName)
	return err
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...int) (err error) {
	var (
		al *alias
	)
	err = pool.RegisterMgoPool(aliasName, dataSource, params...)
	if err != nil {
		DebugLog.Println(err.Error())
		return
	}

	al, err = addAlias(aliasName, driverName)
	if err != nil {
		DebugLog.Println(err.Error())
		return
	}

	al.DataSource = dataSource
	al.DbName = getDatabase(al.DataSource)

	detectTZ(al)

	return
}

// RegisterDriver Register a database driver use specify driver name, this can be definition the driver is which database type.
func RegisterDriver(driverName string, typ DriverType) error {
	if t, ok := drivers[driverName]; !ok {
		drivers[driverName] = typ
	} else {
		if t != typ {
			return fmt.Errorf("driverName `%s` db driver already registered and is other type", driverName)
		}
	}
	return nil
}

// SetDataBaseTZ Change the database default used timezone
func SetDataBaseTZ(aliasName string, tz *time.Location) error {
	if al, ok := dataBaseCache.get(aliasName); ok {
		al.TZ = tz
	} else {
		return fmt.Errorf("DataBase alias name `%s` not registered", aliasName)
	}
	return nil
}

func getDatabase(uri string) (dbName string) {
	cs, err := connstring.Parse(uri)
	if err != nil {
		dbName = "test"
		return
	}
	dbName = cs.Database
	return
}
