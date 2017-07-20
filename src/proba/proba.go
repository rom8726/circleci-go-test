package proba

import "fmt"

type Application struct {
}

func (self *Application) Start() {
	fmt.Println("start!")
}
