package proba

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/pg.v5"
)

type Application struct {
	Database  *pg.DB
	RedisPool *redis.Pool
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
