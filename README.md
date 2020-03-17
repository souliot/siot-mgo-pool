# siot-mgo-pool
A mongodb pool by golang

## uri example
mongodb://yapi:abcd1234@vm:27017/yapi
mongodb://yapi:abcd1234@vm:27017,yapi:abcd1234@vm:27017,yapi:abcd1234@vm:27017/yapi

## mongo Transaction

```golang
  s,err:=db.Begin()
  // Do Something
  if err != nil {
    s.AbortTransaction(todo)
  } else {
    s.CommitTransaction(todo)
  }
```

20200317-工作日志
基于"go.mongodb.org/mongo-driver"的go mongodb驱动，封装数据库操作的orm类库，
目前实现：
1.数据库连接池；
2.orm的（read,ReadOrCreate,Insert,InsertMulti,Update,Delete）方法；
3.querySet的(Filter,Limit,Offset,OrderBy,One,All)方法；
