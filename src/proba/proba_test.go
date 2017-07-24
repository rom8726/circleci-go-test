package proba

import (
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
