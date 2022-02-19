package db

import "xorm.io/xorm"

func Connect() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql","trung:1234@/test1")
	if err != nil {
		panic(err)
	}
	return engine
}
