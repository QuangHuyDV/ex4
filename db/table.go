package db

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"xorm.io/xorm"
)

type Users struct {
	ID string 	 	 `xorm:"pk"`
	Name string
	Birth int64
	Created int64 	 `xorm:"created"`
	UpdatedAt int64 `xorm:"updated_at"`
}

type Point struct {
	UserId string 
	Points int64	
	MaxPoints int64
}

//create
func CreateTable(engine *xorm.Engine, tb interface{}) error {
	_, err := engine.IsTableExist(tb)
	if err != nil {
		log.Println(err)
		return err
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

func InsertUP(engine *xorm.Engine, user interface{}, poi interface{}) error {
	err := InsertTable(engine,"users",user)
	if err != nil {
		return err
	}
	err1 := InsertTable(engine,"point",poi)
	if err1 != nil {
		return err1
	}
	return nil
}

func Insert100(engine *xorm.Engine) error {
	for i := 1; i <= 100; i++ {
		id := "a" + strconv.Itoa(i)
		user := Users{
			ID:  id,
			Name: "trung" + strconv.Itoa(i),
			Created: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		}
		poi := Point{
			UserId: id,
			Points: 10,
			MaxPoints: 10,
		}
		err := InsertUP(engine, user, poi)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

//delete
func DeleteUser(engine *xorm.Engine, tb interface{}, id string) error {
	_, err := engine.Where("i_d = ?", id).Delete(tb)
	if err != nil {
		return err
	}
	return nil
}
func DeletePoint(engine *xorm.Engine, tb interface{}, id string) error {
	_, err := engine.Where("user_id = ?", id).Delete(tb)
	if err != nil {
		return err
	}
	return nil
}

func Delete100(engine *xorm.Engine){
	for i:= 1; i <= 100; i++ {
		id := "a" + strconv.Itoa(i)
		users := Users{}
		point := Point{}
		err := DeleteUser(engine, users, id)
		err1 := DeletePoint(engine, point, id)
		if err != nil || err1 != nil {
			log.Println(err)
		}
	}
}

//update
func UpdateUser(engine *xorm.Engine, data interface{}, id string) error {
	_, err := engine.Table("users").ID(id).Update(data)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePoint(engine *xorm.Engine, data interface{}, id string) error {
	_, err := engine.Table("point").Where("user_id = ?",id).Update(data)
	if err != nil {
		return err
	}
	return nil
}

//read 
func ReadUser(engine *xorm.Engine, id string, tb Users) (Users) {
	_, err := engine.ID(id).Get(&tb)
	if err != nil {
		log.Println(err)
	}
	return tb
}

func ReadPoint(engine *xorm.Engine, id string, tb Point) (Point) {
	_, err := engine.Where("user_id = ?",id).Get(&tb)
	tb = Point{}
	if err != nil {
		log.Println(err)
	}
	return tb
}

//list 
func ListUser(engine *xorm.Engine, tb []Users) ([]Users, int64) {
	count, err := engine.FindAndCount(&tb)
	if err != nil {
		log.Println(err)	
	}
	return tb ,count
}

//scan
func SacnTable(engine *xorm.Engine, id string, tb Users ) (*Users) {
	rows, err := engine.Where("i_d = ?",id).Rows(&tb)
	if err != nil {
		log.Println(err)
	}
	us := new(Users)
	for rows.Next() {
		err = rows.Scan(us)
		if err != nil {
			fmt.Println(err)
		}
	}
	return us
}

// updateBirth
func UpdateBirth(engine *xorm.Engine, id string, dob int64) error {
	session := engine.NewSession()
	defer session.Close()
	er := session.Begin()
	if er != nil {
		session.Rollback()
		return er
	}
	user := Users{
		Birth: dob,
		UpdatedAt: time.Now().UnixMilli(),
	}
	err := UpdateUser(engine,user,id)
	if err != nil {
		session.Rollback()
		return err
	}
	row := ReadPoint(engine,id,Point{})
	uppoint := Point{
		Points: row.Points + int64(10),
	}
	err1 := UpdatePoint(engine, uppoint, id)
	if err1 != nil {
		session.Rollback()
		return err1
	}
	row1 := ReadUser(engine,id, Users{})
	upuser := Users{
		Name: row1.Name + strconv.FormatInt(row1.UpdatedAt,10),
		UpdatedAt: time.Now().UnixMilli(),
	}
	err2 := UpdateUser(engine, upuser, id)
	if err != nil {
		session.Rollback()
		return err2
	}
	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}
