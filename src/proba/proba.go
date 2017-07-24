package proba

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/pg.v5"
	"strconv"
	"strings"
)

type Application struct {
	Database        *pg.DB
	RedisPool       *redis.Pool
	Couchbase       *gocb.Bucket
	AerospikeClient *aerospike.Client
}

func NewApplication() (application Application) {
	var err error
	application.Database, err = NewPostgreSqlClient("127.0.0.1", 5432, "circleci", "", "circleci-go-test", 2)
	if err != nil {
		panic(err)
	}

	application.RedisPool, err = NewRedisPool("127.0.0.1", 6379, 2)
	if err != nil {
		panic(err)
	}

	_, application.Couchbase, err = NewCouchbaseClient("127.0.0.1", "test", "")
	if err != nil {
		panic(err)
	}

	application.AerospikeClient, err = NewAerospikeClient([]string{"127.0.0.1:3000"})
	if err != nil {
		panic(err)
	}

	return
}

func (self *Application) Start() {
	fmt.Println("start!")
}

func (self *Application) Close() {
	if self.Database != nil {
		self.Database.Close()
	}
	if self.RedisPool != nil {
		self.RedisPool.Close()
	}
	if self.Couchbase != nil {
		self.Couchbase.Close()
	}
	if self.AerospikeClient != nil {
		self.AerospikeClient.Close()
	}
}

func (self *Application) SomeFunc() int {
	return 3
}

func (self *Application) RedisFunc() error {
	conn := self.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", "queue", "start")
	return err
}

func (self *Application) CouchbaseFunc() error {
	_, err := self.Couchbase.Upsert("test-key", "test-value", 60)
	return err
}

func (self *Application) AerospikeFunc() error {
	key, err := aerospike.NewKey("test", "test", "test-key")
	if err != nil {
		return err
	}
	operations := []*aerospike.Operation{
		aerospike.PutOp(aerospike.NewBin("bin1", 1)),
		aerospike.PutOp(aerospike.NewBin("bin2", 2)),
		aerospike.AddOp(aerospike.NewBin("metric", 12)),
	}
	_, err = self.AerospikeClient.Operate(nil, key, operations...)
	return err
}

// NewPostgreSqlClient initializes connection to database for pg.DB (models)
func NewPostgreSqlClient(host string, port int, username string, password string, database string, max_conns int) (conn *pg.DB, err error) {
	conn = pg.Connect(&pg.Options{
		PoolSize:   max_conns + 1,
		User:       username,
		Password:   password,
		Database:   database,
		Addr:       fmt.Sprintf("%s:%d", host, port),
		MaxRetries: 3,
	})

	// ping
	_, err = conn.Exec("SELECT 'ping'")
	if err != nil {
		return nil, err
	}

	return
}

// NewRedisPool initializes pool of connections for Redis
func NewRedisPool(host string, port int, max_conns int) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     max_conns,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), redis.DialDatabase(0))
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	// ping
	conn := pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("PING"))
	if err != nil {
		return nil, err
	}
	if reply != "PONG" {
		return nil, err
	}

	return
}

// NewCouchbaseClient create connection to Couchbase bucket
func NewCouchbaseClient(host string, bucket_name string, password string) (cluster *gocb.Cluster, bucket *gocb.Bucket, err error) {
	cluster, err = gocb.Connect(fmt.Sprint("couchbase://", host))
	if err != nil {
		return
	}

	bucket, err = cluster.OpenBucket(bucket_name, password)
	return
}

// NewAerospikeClient returns new Aerospike client
func NewAerospikeClient(nodes []string) (as_client *aerospike.Client, err error) {
	hosts := []*aerospike.Host{}

	for _, server := range nodes {
		parts := strings.SplitN(server, ":", 2)
		port := 3000
		if len(parts) == 2 {
			if port, err = strconv.Atoi(parts[1]); err != nil {
				return
			}
		}
		hosts = append(hosts, &aerospike.Host{Name: parts[0], Port: port})
	}
	as_client, err = aerospike.NewClientWithPolicyAndHost(nil, hosts...)
	return
}
