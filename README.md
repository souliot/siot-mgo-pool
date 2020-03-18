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

[
map[
key:map[_id:1]
name:_id_
ns:yapi.log v:2
]
map[
key:map[type:1 typeid:1]
name:typeid_1_type_1 ns:yapi.log v:2]
map[key:map[uid:1] name:uid_1 ns:yapi.log v:2]
]
