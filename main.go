package main

import (
	"ex4/db"
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

func insertUser(engine *xorm.Engine, id string, name string, dob int64) error {
	user := Users{
		ID: id,
		Name: name,
		Birth: dob,
		Created: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	err := db.InsertTable(engine,"users",user)
	if err != nil {
		return err
	} else {
		poi := Point{
			UserId: id,
			Points: 10,
			MaxPoints: 10,
		}
		err := db.InsertTable(engine,"point",poi)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateUser(engine *xorm.Engine, id string) error {
	now := time.Now().UnixMilli()
	user := Users{
		UpdatedAt: now,
	}
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

func updateBirth(engine *xorm.Engine, id string, dob int64) error {
	session := engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	user := Users{
		Birth: dob,
		UpdatedAt: time.Now().UnixMilli(),
	}
	err := db.UpdateTable(engine,"users",user,id)
	if err != nil {
		return err
	} else {
		row := readPoint(engine,id)
		uppoint := Point{
			Points: row.Points + 10,
		}
		_, err := engine.Table("point").Where("user_id = ?",id).Update(&uppoint)
		if err != nil {
			return err
		} else {
			row1 := readUser(engine,id)
			upuser := Users{
				Name: row1.Name + strconv.FormatInt(row1.UpdatedAt,10),
				UpdatedAt: time.Now().UnixMilli(),
			}
			err := db.UpdateTable(engine, "users", upuser, id)
			if err != nil {
				return err
			}
		}
	}
	err = session.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func insert100(engine *xorm.Engine) error {
	for i := 0; i < 100; i++ {
		id := "a" + strconv.Itoa(i)
		name := "trung" + strconv.Itoa(i)
		insertUser(engine, id, name, 1234)
	}
	return nil
}

func delete100(engine *xorm.Engine){
	for i:= 0; i <100; i++ {
		id := "a" + strconv.Itoa(i)
		users := Users{}
		point := Point{}
		err := db.DeleteTable(engine, users, id)
		err1 := db.DeleteTable(engine, point, id)
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
		panic(err)
	}
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		defer wg.Done()
		id := "a" + strconv.Itoa(i)
		us := readUser(engine, id)
		name <- us.Name
		id1 <- us.ID
		go log.Printf("$(%v) - $(%v) - $(%v)", i, <- id1, <-name)
	}
	close(name)
	close(id1)  
}

func main(){
	engine := db.Connect()
	defer engine.Close()
	// engine.ShowSQL(true)
	//ex1: 

	// 1. Viết hàm: Chỉ tạo db, và tạo model(struct) ánh xạ struct thành table (CreateTable, Sync2)

	// table.CreateTable(engine, Users{})
	// table.CreateTable(engine, Point{})
	
	// table.DropTable(engine,Users{})
	// table.DropTable(engine,Point{})

	//2. Viết hàm: insert và update user, viết hàm list user hoặc đọc user theo id(4 hàm)
	//3. Viết hàm: sau khi tạo user thì insert user_id vào user_point với số điểm 10.

	// err := insertUser(engine, "s3", "minh", 1234)
	// if err != nil {
	// 	log.Println(err)
	// }

	// updateUser(engine, "s2")

	// db.ListTable(engine,"users")

	// u := readUser(engine, "a10")
	// log.Printf("User: \nid: %v\nName: %v\nBirth: %v\n ",u.ID,u.Name,u.Birth)

	//ex2: tạo 1 transaction khi update birth thành công thì cộng 10 điểm vào point sau đó sửa lại name thành $name + "updated " nếu 1 quá trình fail thì rollback, xong commit (CreateSesson)

	// err := updateBirth(engine, "s2", 9999)
	// if err != nil {
	// 	log.Println(err)
	// }

	//ex3: insert 100 bản ghi vào user: sau đó viết 1 workerpool scantableuser lấy ra tên của các user inra màn hình (Dùng scan theo row) dùng 2 worker và thiết lập bộ đếm ${counter} - ${id} - ${name}

	// insert100(engine)
	// scantableuser(engine)

}