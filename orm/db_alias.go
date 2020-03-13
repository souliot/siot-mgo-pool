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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/souliot/siot-mgo-pool/pool"
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
	*sync.RWMutex
	MDB *mongo.Database
}

var _ dbQuerier = new(DB)

func (d *DB) Begin() (s mongo.Session, err error) {
	s, err = d.MDB.Client().StartSession()
	if err != nil {
		return
	}
	defer s.EndSession(todo)

	//开始事务
	err = s.StartTransaction()

	return
}

func (d *DB) GetDB() (s *mongo.Database) {
	return d.MDB
}

type alias struct {
	Name         string
	Driver       DriverType
	DriverName   string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DB           *DB
	DbBaser      dbBaser
	TZ           *time.Location
	Engine       string
}

func detectTZ(al *alias) {
	// orm timezone system match database
	// default use Local
	al.TZ = DefaultTimeLoc
}

func addAliasWthDB(aliasName, driverName string, db *mongo.Database) (*alias, error) {
	al := new(alias)
	al.Name = aliasName
	al.DriverName = driverName
	al.DB = &DB{
		RWMutex: new(sync.RWMutex),
		MDB:     db,
	}

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
func AddAliasWthDB(aliasName, driverName string, db *mongo.Database) error {
	_, err := addAliasWthDB(aliasName, driverName, db)
	return err
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...int) (err error) {
	var (
		al *alias
	)
	pool.RegisterMgoPool(aliasName, dataSource, params)
	client, err := pool.GetMgoClient(aliasName)
	if err != nil {
		DebugLog.Println(err.Error())
		return
	}

	dbNames, err := client.ListDatabaseNames(todo, bson.D{})
	if err != nil {
		DebugLog.Println(err.Error())
		return
	}
	db := client.Database(dbNames[0])

	al, err = addAliasWthDB(aliasName, driverName, db)
	if err != nil {
		DebugLog.Println(err.Error())
		return
	}

	al.DataSource = dataSource

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

// GetDB Get *sql.DB from registered database by db alias name.
// Use "default" as alias name if you not set.
func GetDB(aliasNames ...string) (*mongo.Database, error) {
	var name string
	if len(aliasNames) > 0 {
		name = aliasNames[0]
	} else {
		name = "default"
	}
	al, ok := dataBaseCache.get(name)
	if ok {
		return al.DB.MDB, nil
	}
	return nil, fmt.Errorf("DataBase of alias name `%s` not found", name)
}
