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
