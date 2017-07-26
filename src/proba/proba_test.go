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
		So(app.SomeFunc(), ShouldEqual, 3)
	})
}

func TestApplication_PostgreFunc(t *testing.T) {
	Convey("PostgreFunc() should work correctly", t, func() {
		app := NewApplication()
		So(app.PostgreFunc(), ShouldBeNil)
	})
}

func TestApplication_RedisFunc(t *testing.T) {
	Convey("RedisFunc() should work correctly", t, func() {
		app := NewApplication()

		So(app.RedisFunc(), ShouldBeNil)

		pool, err := NewRedisPool()
		So(err, ShouldBeNil)
		defer pool.Close()
		conn := pool.Get()
		defer conn.Close()
		res, err := redis.String(conn.Do("RPOP", "queue"))
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "start")
	})
}

func TestApplication_AerospikeFunc(t *testing.T) {
	Convey("AerospikeFunc() should work correctly", t, func() {
		app := NewApplication()

		So(app.AerospikeFunc(), ShouldBeNil)

		as_client, err := NewAerospikeClient()
		So(err, ShouldBeNil)
		defer as_client.Close()
		key, err := aerospike.NewKey("test", "test", "test-key")
		So(err, ShouldBeNil)
		rec, err := as_client.Get(nil, key, "bin1", "bin2", "metric")
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

		So(app.CouchbaseFunc(), ShouldBeNil)

		_, bucket, err := NewCouchbaseClient()
		So(err, ShouldBeNil)
		defer bucket.Close()

		res := ""
		_, err = bucket.Get("test-key", &res)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "test-value")
	})
}

func TestApplication_KafkaProducerFunc(t *testing.T) {
	Convey("KafkaProducerFunc() should work correctly", t, func() {
		app := NewApplication()

		So(app.KafkaProducerFunc(), ShouldBeNil)
	})
}

func TestApplication_KafkaConsumerFunc(t *testing.T) {
	Convey("KafkaConsumerFunc() should work correctly", t, func() {
		app := NewApplication()

		So(app.KafkaConsumerFunc(), ShouldBeNil)
	})
}
