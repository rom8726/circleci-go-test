package proba

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApplication_SomeFunc(t *testing.T) {
	Convey("SomeFunc() should work correctly", t, func() {
		app := NewApplication()
		defer app.Close()
		So(app.SomeFunc(), ShouldEqual, 3)
	})
}

func TestApplication_RedisFunc(t *testing.T) {
	Convey("RedisFunc() should work correctly", t, func() {
		app := NewApplication()
		defer app.Close()

		So(app.RedisFunc(), ShouldBeNil)

		conn := app.RedisPool.Get()
		defer conn.Close()
		res, err := redis.String(conn.Do("RPOP", "queue"))
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "start")
	})
}

func TestApplication_AerospikeFunc(t *testing.T) {
	Convey("AerospikeFunc() should work correctly", t, func() {
		app := NewApplication()
		defer app.Close()

		So(app.AerospikeFunc(), ShouldBeNil)

		key, err := aerospike.NewKey("test", "test", "test-key")
		So(err, ShouldBeNil)
		rec, err := app.AerospikeClient.Get(nil, key, "bin1", "bin2", "metric")
		So(err, ShouldBeNil)
		So(rec.Bins["bin1"], ShouldNotBeNil)
		So(rec.Bins["bin2"], ShouldNotBeNil)
		So(rec.Bins["metric"], ShouldNotBeNil)
		So(rec.Bins["bin1"].(int), ShouldEqual, 1)
		So(rec.Bins["bin2"].(int), ShouldEqual, 2)
		So(rec.Bins["metric"].(int), ShouldEqual, 12)
	})
}

func TestApplication_CouchbaseFunc(t *testing.T) {
	Convey("CouchbaseFunc() should work correctly", t, func() {
		app := NewApplication()
		defer app.Close()

		So(app.CouchbaseFunc(), ShouldBeNil)

		res := ""
		_, err := app.Couchbase.Get("test-key", &res)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "test-value")
	})
}
