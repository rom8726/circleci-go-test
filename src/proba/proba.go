package proba

import (
	"fmt"
	"gopkg.in/pg.v5"
)

type Application struct {
	Database *pg.DB
}

func NewApplication() (application Application) {
	var err error
	application.Database, err = NewPostgreSqlClient("127.0.0.1", 5432, "circleci", "", "proba", 2)
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
}

func (self *Application) SomeFunc() int {
	return 3
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
