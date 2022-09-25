package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type _post struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	User    string `json:"user"`
} //推文结构
type _user struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Password string `json:"password"`
} //用户结构
type _comment struct {
	Id      int64  `json:"id"`
	Comment string `json:"comment"`
	Pid     int64  `json:"pid"`
} //评论结构
var postNums int
var userNums int
var commentNums int

func main() {
	InitDB()           //初始化数据库
	r := setupRouter() //设置路由
	r.Run()            //启动服务
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/post/addPost", AddPost)
	r.GET("/api/post/getPost", GetPost)
	r.PUT("/api/post/putPost", PutPost)
	r.DELETE("/api/post/deletePost", DeletePost)
	r.POST("/api/user/register", Register)
	r.POST("/api/user/login", Login)
	r.POST("/api/comment/addComment", AddComment)
	r.GET("/api/comment/getComment", GetComment)
	r.DELETE("/api/comment/deleteComment", DeleteComment)
	return r
}
func AddPost(c *gin.Context) { //添加推文
	var post _post
	c.ShouldBindJSON(&post) //绑定数据到自定义结构体
	post.Time = time.Now().Format("2006/1/02/ 15:04")
	//存入数据库
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	add, err := db.Prepare("INSERT INTO posts(content,name,time,user) VALUES (?,?,?,?)")
	CheckErr(err)
	add.Exec(post.Content, post.Name, post.Time, post.User)
	postNums++
	db.Close()
	{
		//return response
	}
}
func GetPost(c *gin.Context) { //获取推文

}
func PutPost(c *gin.Context) { //更新推文
	var post _post
	c.ShouldBindJSON(&post)
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	put, err := db.Prepare("UPDATE posts set content=? where id=?")
	CheckErr(err)
	put.Exec(post.Content, post.Id)
	db.Close()
	{
		//return response
	}
}
func DeletePost(c *gin.Context) { //删除推文
	id := c.Query("Id")
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	delete, err := db.Prepare("DELETE FROM posts where id=?")
	delete.Exec(id)
	postNums--
	{
		//return response
	}
	db.Close()
}

func Register(c *gin.Context) { //注册
	var re int
	var user _user
	c.ShouldBindJSON(&user)
	var userFromDB string
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	rows, err := db.Query("SELECT username FROM users")
	//检测名称重复
	for rows.Next() {
		rows.Scan(&userFromDB)
		if user.Username == userFromDB {
			re = 1
			{
				//return response
			}
		}
	}
	if re != 1 {
		add, err := db.Prepare("INSERT INTO users(username, password) VALUES (?,?)")
		CheckErr(err)
		add.Exec(user.Username, user.Password)
		{
			//return response
		}
	}
	db.Close()
}
func Login(c *gin.Context) { //登录

}

func AddComment(c *gin.Context) { //添加评论
	var comment _comment
	c.ShouldBindJSON(&comment)
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	add, err := db.Prepare("INSERT INTO comments(comment, pid) VALUES (?,?)")
	add.Exec(comment.Comment, comment.Pid)
	{
		//return response
	}
	db.Close()
}
func GetComment(c *gin.Context) { //获取评论
	var comments []_comment
	var pid int64
	var commentFromDB _comment
	c.Bind(pid)
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	rows, err := db.Query("SELECT pid FROM comments")
	for rows.Next() {
		rows.Scan(&commentFromDB)
		if pid == commentFromDB.Pid {
			comments = append(comments, commentFromDB)
		}
	}
	if comments == nil {
		//return response
	} else {
		//return data
	}
	db.Close()
}
func DeleteComment(c *gin.Context) { //删除评论
	var id int64
	c.Bind(&id)
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	delete, err := db.Prepare("DELETE FROM comments WHERE ID=?")
	CheckErr(err)
	delete.Exec(id)
	commentNums--
	{
		//return response
	}
	db.Close()
}
func InitDB() {
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	sqlTablePost := `CREATE TABLE IF NOT EXISTS "posts"(
	    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
	    "content" VARCHAR(1024) NULL,
	    "name" VARCHAR(20) NULL,
	    "time" VARCHAR(50) NULL,
	    "user" VARCHAR(100) NULL
	)`
	db.Exec(sqlTablePost)
	sqlTableUser := `
	CREATE TABLE IF NOT EXISTS "users"(
	    "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
	    "username" VARCHAR(100) NULL,
	    "password" VARCHAR(100) NULL
	)`
	db.Exec(sqlTableUser)
	sqlTableComment := `
	CREATE TABLE IF NOT EXISTS "comments"(
	    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
	    "comment" VARCHAR(100) NULL,
	    "pid" INTEGER NULL
	)`
	db.Exec(sqlTableComment)
	db.Close()
}
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
