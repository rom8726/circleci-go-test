package proba

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aerospike/aerospike-client-go"
	"github.com/garyburd/redigo/redis"
	"github.com/wvanbergen/kafka/consumergroup"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/pg.v5"
	"time"
)

const (
	PG_USERNAME = "circleci"
	PG_PASSWORD = ""
	PG_DATABASE = "circleci-go-test"

	COUCHBASE_BUCKET = "test"

	KAFKA_TOPIC = "my-topic"
)

type Application struct {
}

func NewApplication() (application Application) {
	return
}

func (self *Application) Start() {
	fmt.Println("start!")
}

func (self *Application) SomeFunc() int {
	return 3
}

func (self *Application) PostgreFunc() error {
	db, err := NewPostgreSqlClient()
	if err != nil {
		return err
	}
	defer db.Close()

	// ping
	_, err = db.Exec("SELECT 'ping'")
	return err
}

func (self *Application) RedisFunc() error {
	pool, err := NewRedisPool()
	if err != nil {
		return err
	}
	defer pool.Close()

	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("LPUSH", "queue", "start")
	return err
}

func (self *Application) CouchbaseFunc() error {
	_, bucket, err := NewCouchbaseClient()
	if err != nil {
		return err
	}
	defer bucket.Close()

	_, err = bucket.Upsert("test-key", "test-value", 60)
	return err
}

func (self *Application) AerospikeFunc() error {
	as_client, err := NewAerospikeClient()
	if err != nil {
		return err
	}
	defer as_client.Close()

	key, err := aerospike.NewKey("test", "test", "test-key")
	if err != nil {
		return err
	}
	operations := []*aerospike.Operation{
		aerospike.PutOp(aerospike.NewBin("bin1", 1)),
		aerospike.PutOp(aerospike.NewBin("bin2", 2)),
		aerospike.AddOp(aerospike.NewBin("metric", 12)),
	}
	_, err = as_client.Operate(nil, key, operations...)
	return err
}

func (self *Application) KafkaProducerFunc() error {
	kafka_config := sarama.NewConfig()
	kafka_config.ClientID = "KAFKA_CLIENT_ID"
	kafka_config.Producer.Return.Successes = true
	kafka_config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, kafka_config)
	if err != nil {
		return err
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: KAFKA_TOPIC,
		Value: sarama.StringEncoder("test message"),
	}
	_, _, err = producer.SendMessage(msg)
	return err
}

func (self *Application) KafkaConsumerFunc() error {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = time.Minute
	config.Offsets.CommitInterval = time.Second
	config.Offsets.ResetOffsets = false

	consumer, err := consumergroup.JoinConsumerGroup(
		"group_id",
		[]string{KAFKA_TOPIC},
		[]string{"127.0.0.1:2181"},
		config,
	)
	if err != nil {
		return err
	}
	if consumer.Closed() {
		return fmt.Errorf("Kafka consumer is closed!")
	}

	var registered bool
	registered, err = consumer.InstanceRegistered()
	if err != nil {
		return err
	}
	if !registered {
		return fmt.Errorf("Kafka consumer is not registered!")
	}

	defer consumer.Close()
	select {
	case message := <-consumer.Messages():
		fmt.Printf("Kafka message: %+v\n", message)
		break
	default:
		break
	}

	return nil
}

// NewPostgreSqlClient initializes connection to database for pg.DB (models)
func NewPostgreSqlClient() (conn *pg.DB, err error) {
	conn = pg.Connect(&pg.Options{
		PoolSize:   2,
		User:       PG_USERNAME,
		Password:   PG_PASSWORD,
		Database:   PG_DATABASE,
		Addr:       fmt.Sprintf("127.0.0.1:5432"),
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
func NewRedisPool() (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("127.0.0.1:6379"), redis.DialDatabase(0))
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
func NewCouchbaseClient() (cluster *gocb.Cluster, bucket *gocb.Bucket, err error) {
	cluster, err = gocb.Connect("couchbase://127.0.0.1")
	if err != nil {
		return
	}

	bucket, err = cluster.OpenBucket(COUCHBASE_BUCKET, "")
	return
}

// NewAerospikeClient returns new Aerospike client
func NewAerospikeClient() (as_client *aerospike.Client, err error) {
	//as_host_out, err := exec.Command("sh", "-c", "docker inspect -f '{{.NetworkSettings.IPAddress }}' aerospike").Output()
	//if err != nil {
	//	return nil, err
	//}
	//as_host := string(as_host_out)
	//as_host = strings.Replace(as_host, "\n", "", -1)
	//fmt.Println(fmt.Sprint("Aerospike host: ", string(as_host)))

	as_client, err = aerospike.NewClientWithPolicyAndHost(nil, &aerospike.Host{
		Name: "172.29.0.3",
		Port: 3000,
	})
	return
}
