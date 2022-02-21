package main

import (
	"ex4/db"
	"fmt"
	"log"
	"strconv"
	"sync"

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


func scantableuser(engine *xorm.Engine) {
 	var wg sync.WaitGroup
	name := make(chan string,3)
	id1 := make(chan string,3)
	u := Users{}
	count,_ := engine.Count(&u)
	for i := 1; i <= int(count); i++ {
		wg.Add(1)
		defer wg.Done()
		id := "a" + strconv.Itoa(i)
		us := db.SacnTable(engine, id, db.Users{})
		name <- us.Name
		id1 <- us.ID
		go worker(i, id1, name)
	}
	close(name)
	close(id1)  
}

func worker(i int,id1 <- chan string, name <- chan string)  {
	fmt.Printf("%v - %v - %v\n", i, <- id1, <-name)
}

func main(){
	engine, err := db.Connect()
	if err != nil {
		log.Println(err)
	}
	defer engine.Close()
	// engine.ShowSQL(true)
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
	// id := "d2"
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

	// err = db.InsertUP(engine, user, poi)
	// if err != nil {
	// 	log.Println(err)
	// }

		//Update
	// id:= "d2"
	// user := Users{
	// 	ID: id,
	// 	Name: "Tran",
	// 	UpdatedAt: time.Now().UnixMilli(),
	// }
	// err = db.UpdateUser(engine, user, id)
	// if err != nil {
	// 	log.Println(err)
	// }

		//List
	// u, count := db.ListUser(engine, []db.Users{})
	// for i := 0; i < int(count); i++ {
	// 	fmt.Printf("User%v: id: %v, Name: %v, Birth: %v, Created: %v, Updated_at: %v\n", i, u[i].ID, u[i].Name,u[i].Birth, u[i].Created, u[i].UpdatedAt)
	// }

		//Read
	// us := db.ReadUser(engine, "d1", db.Users{})
	// fmt.Printf("User: %v, name: %v, birth: %v, created: %v, updated_at: %v",us.ID,us.Name,us.Birth,us.Created,us.UpdatedAt)

	//ex2: tạo 1 transaction khi update birth thành công thì cộng 10 điểm vào point sau đó sửa lại name thành $name + "updated " nếu 1 quá trình fail thì rollback, xong commit (CreateSesson)

	// err1 := db.UpdateBirth(engine, "s3", 221112)
	// if err1 != nil {
	// 	log.Println(err1)
	// }

	//ex3: insert 100 bản ghi vào user: sau đó viết 1 workerpool scantableuser lấy ra tên của các user inra màn hình (Dùng scan theo row) dùng 2 worker và thiết lập bộ đếm ${counter} - ${id} - ${name}

	// db.Insert100(engine)
	// db.Delete100(engine)
	// scantableuser(engine)

}
