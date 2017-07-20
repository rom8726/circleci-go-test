package proba

import "fmt"

type Application struct {
}

func (self *Application) Start() {
	fmt.Println("start!")
}

func (self *Application) SomeFunc() int {
	return 3
}
