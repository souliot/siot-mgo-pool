package orm

import (
	"testing"

	"github.com/souliot/siot-mgo-pool/pool"
)

func init() {
	pool.RegisterMgoPool("default", "mongodb://yapi:abcd1234@vm:27017/yapi")
}
func TestDB(t *testing.T) {

}
