package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	_"github.com/go-sql-driver/mysql"
)

type Todo struct {
	gorm.Model
	Text string
	Status string
}


//DBMigrate
func dbInit(){
	db, err := gorm.Open("mysql","root@/sample?charset=utf8&parseTime=True&loc=Local")

	defer db.Close()


	if err != nil {
		panic("Failed to connect database")
	}

	if err == nil {
		print("Success!!")
	}

	db.AutoMigrate(&Todo{})
}

//DB追加
func dbInsert(text string, status string) {
	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbInserteErr")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

//DB更新
func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbUpdateErr")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

//DB削除
func dbDelete(id int) {
	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbDeleteErr")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

//DB全取得
func dbGetAll() []Todo {
	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbGetAllErr")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

//DB一つ取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbGetOneErr")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}


func main() {


	router := gin.Default()

	//静的ファイルを有効にする
	router.Static("/assets", "./assets")
	//htmlを有効に
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		todos := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(302, "/")
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(302, "/")
	})

	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")

	})

	router.Run()
}




