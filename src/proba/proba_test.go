package proba

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApplication_SomeFunc(t *testing.T) {
	Convey("SomeFunc() should work correctly", t, func() {
		app := NewApplication()
		app.Close()
		So(app.SomeFunc(), ShouldEqual, 3)

		conn := app.RedisPool.Get()
		defer conn.Close()
		res, err := redis.String(conn.Do("RPOP", "queue"))
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "start")
	})
}
