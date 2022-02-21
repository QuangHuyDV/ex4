package main

import (
	"ex4/db"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func insertUser(engine *xorm.Engine, user interface{}, poi interface{}) error {
	err := db.InsertTable(engine,"users",user)
	if err != nil {
		return err
	}
	err1 := db.InsertTable(engine,"point",poi)
	if err1 != nil {
		return err1
	}
	return nil
}

func updateUser(engine *xorm.Engine, user interface{}, id string) error {
	err := db.UpdateTable(engine, "users", user, id)
	if err != nil {
		return err
	}
	return nil
}

func readUser(engine *xorm.Engine, id string) *Users {
	u := new(Users)
	rows, err := engine.Where("i_d = ?",id).Rows(u)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(u)
	}
	return u
}

func readPoint(engine *xorm.Engine, id string) *Point {
	p := new(Point)
	rows, err := engine.Where("user_id = ?",id).Rows(p)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(p)
	}
	return p
}

func ListUser(engine *xorm.Engine) {
	u := []Users{}
	err := engine.Find(&u)
	if err != nil {
		log.Println(err)
	}
	us:= Users{}
	count,err := engine.Count(&us)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < int(count); i++ {
		fmt.Printf("User%v: id: %v, Name: %v, Birth: %v, Created: %v, Updated_at: %v\n", i, u[i].ID, u[i].Name,u[i].Birth, u[i].Created, u[i].UpdatedAt)
	}
}

func updateBirth(engine *xorm.Engine, id string, dob int64) error {
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
	_,err := session.Table("users").ID(id).Update(&user)
	if err != nil {
		session.Rollback()
		return err
	}
	row := readPoint(engine,id)
	uppoint := Point{
		Points: row.Points + int64(10),
	}
	_, err1 := session.Table("point").Where("user_id = ?",id).Update(&uppoint)
	if err1 != nil {
		session.Rollback()
		return err1
	}
	row1 := readUser(engine,id)
	upuser := Users{
		Name: row1.Name + strconv.FormatInt(row1.UpdatedAt,10),
		UpdatedAt: time.Now().UnixMilli(),
	}
	err2 := db.UpdateTable(engine, "users", upuser, id)
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

func insert100(engine *xorm.Engine) error {
	for i := 0; i < 100; i++ {
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
		insertUser(engine, user, poi)
	}
	return nil
}

func delete100(engine *xorm.Engine){
	for i:= 0; i <100; i++ {
		id := "a" + strconv.Itoa(i)
		users := Users{}
		point := Point{}
		err := db.DeleteTable(engine, users, id)
		err1 := db.DeleteTable1(engine, point, id)
		if err != nil || err1 != nil {
			log.Println(err, err1)
		}
	}
}

var wg sync.WaitGroup

func scantableuser(engine *xorm.Engine) {
	u := Users{}
	name := make(chan string,3)
	id1 := make(chan string,3)
	count,err := engine.Count(&u)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		defer wg.Done()
		id := "a" + strconv.Itoa(i)
		us := readUser(engine, id)
		name <- us.Name
		id1 <- us.ID
		go worker(i,name, id1)
	}
	close(name)
	close(id1)  
}

func worker(i int,id1 <- chan string, name <- chan string)  {
	log.Printf("%v - %v - %v", i, <- id1, <-name)
}

func main(){
	engine, err := db.Connect()
	if err != nil {
		log.Println(err)
	}
	defer engine.Close()
	engine.ShowSQL(true)
	//ex1: 

	// 1. Viết hàm: Chỉ tạo db, và tạo model(struct) ánh xạ struct thành table (CreateTable, Sync2)

	// err1 := db.CreateTable(engine, Users{})
	// err1 := db.CreateTable(engine, Point{})
	
	//  err1 := db.DropTable(engine,Users{})
	//  err1 := db.DropTable(engine,Point{})
	// if err1 != nil {
	// 	log.Println(err1)
	// }

	//2. Viết hàm: insert và update user, viết hàm list user hoặc đọc user theo id(4 hàm)
	//3. Viết hàm: sau khi tạo user thì insert user_id vào user_point với số điểm 10.

		//insert
	// id := "d1"
	// user := Users{
	// 	ID: id,
	// 	Name: "Ha",
	// 	Birth: 1234,
	// 	Created: time.Now().UnixMilli(),
	// 	UpdatedAt: time.Now().UnixMilli(),
	// }
	// poi := Point{
	// 	UserId: id,
	// 	Points: 10,
	// 	MaxPoints: 10,
	// }

	// insertUser(engine, user, poi)

		//Update
	// id:= "s2"
	// user := Users{
	// 	ID: id,
	// 	Name: "aaa",
	// 	UpdatedAt: time.Now().UnixMilli(),
	// }
	// updateUser(engine, user, id)

		//List
	// ListUser(engine)

		//Read
	// u := readUser(engine, "a10")
	// log.Printf("User: \nid: %v\nName: %v\nBirth: %v\n ",u.ID,u.Name,u.Birth)

	//ex2: tạo 1 transaction khi update birth thành công thì cộng 10 điểm vào point sau đó sửa lại name thành $name + "updated " nếu 1 quá trình fail thì rollback, xong commit (CreateSesson)

	// err1 := updateBirth(engine, "s2", 222)
	// if err1 != nil {
	// 	log.Println(err1)
	// }

	//ex3: insert 100 bản ghi vào user: sau đó viết 1 workerpool scantableuser lấy ra tên của các user inra màn hình (Dùng scan theo row) dùng 2 worker và thiết lập bộ đếm ${counter} - ${id} - ${name}

	// insert100(engine)
	// delete100(engine)
	// scantableuser(engine)

}
