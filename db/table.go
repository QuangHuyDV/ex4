package db

import (
	"log"

	"xorm.io/xorm"
)

//create
func CreateTable(engine *xorm.Engine, tb interface{}) {
	ts, err := engine.IsTableExist(tb)
	if err != nil {
		log.Println(err)
	} else if ts {
		log.Println("Bang da ton tai")
	} else {
		err := engine.Sync2(tb)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Them thanh cong")
		}
	}
}

//drop
func DropTable(engine *xorm.Engine, tb interface{}) error {
	err := engine.DropTables(tb)
	if err != nil {
		return err
	} 
	return err
}

//insert
func InsertTable(engine *xorm.Engine, tb string, data interface{}) error {
	_, err := engine.Table(tb).Insert(data)
	if err != nil {
		return err
	}
	return err
}

func DeleteTable(engine * xorm.Engine, tb interface{}, id string) error {
	_, err := engine.Where("i_d = ?",id).Delete(tb)
	if err != nil {
		return err
	}
	return nil
}

//update
func UpdateTable(engine *xorm.Engine, tb string, data interface{}, id string) error {
	_, err := engine.Table(tb).ID(id).Update(&data)
	if err != nil {
		return err
	} 
	return err
}

//list
func ListTable(engine *xorm.Engine, tb string) {
	count, err := engine.Table(tb).Count()
	if err != nil {
		log.Println(err)
	}
	has, err := engine.Table(tb).QueryString()
	if err != nil {
		panic(err)
	}
	for i:= 0; i < int(count); i++ {
		log.Printf("User%v: %v\n",i,has[i])
	}
}
