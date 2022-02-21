package db

import (
	"log"

	"xorm.io/xorm"
)

//create
func CreateTable(engine *xorm.Engine, tb interface{}) error {
	ts, err := engine.IsTableExist(tb)
	if err != nil {
		log.Println(err)
		return err

	}
	if ts {
		return nil
	}
	err = engine.Sync2(tb)
	if err != nil {
		log.Println(err)
		return err

	}
	return nil
}

//drop
func DropTable(engine *xorm.Engine, tb interface{}) error {
	err := engine.DropTables(tb)
	if err != nil {
		return err
	}
	return nil
}

//insert
func InsertTable(engine *xorm.Engine, tb string, data interface{}) error {
	_, err := engine.Table(tb).Insert(data)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTable(engine *xorm.Engine, tb interface{}, id string) error {
	_, err := engine.Where("i_d = ?", id).Delete(tb)
	if err != nil {
		return err
	}
	return nil
}
func DeleteTable1(engine *xorm.Engine, tb interface{}, id string) error {
	_, err := engine.Where("user_id = ?", id).Delete(tb)
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
	return nil
}

func UpdateCol(engine * xorm.Engine, tb string, data interface{}, id string) error {
	_, err := engine.Cols("Birth","Updated_at").Table(tb).ID(id).Update(&data)
	if err != nil {
		return err
	}
	return nil
}

