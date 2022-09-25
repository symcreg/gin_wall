package main

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
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
type _userClaims struct {
	username string
	//jwt提供的标准claims
	jwt.StandardClaims
} //jwt struct
var (
	secret     = []byte("gin_wall") //SecretKey
	EffectTime = 2 * time.Hour      //token有效时间
)
var (
	postNums    int //post数量
	userNums    int //user数量
	commentNums int //comment数量
)

func main() {
	InitDB()           //初始化数据库
	r := setupRouter() //设置路由
	r.Run()            //启动服务
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(JwtVerify)     //token中间件
	r.Group("/api/post") //推文路由组
	{
		r.POST("/addPost", AddPost)
		r.GET("/getPost", GetPost)
		r.PUT("/putPost", PutPost)
		r.DELETE("/api/post/deletePost", DeletePost)
	}

	r.Group("/api/user") //用户路由组
	{
		r.POST("/register", Register)
		r.POST("/login", Login)
		r.POST("/auth", auth)
	}
	r.Group("/api/comment") //评论路由组
	{
		r.POST("/api/comment/addComment", AddComment)
		r.GET("/api/comment/getComment", GetComment)
		r.DELETE("/api/comment/deleteComment", DeleteComment)
	}
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
	var user _user
	c.ShouldBindJSON(&user)
	var userFromDB _user
	var exist int = 1
	db, err := sql.Open("sqlite3", "wall.db")
	CheckErr(err)
	rows, err := db.Query("SELECT * FROM users")
	for rows.Next() {
		rows.Scan(&userFromDB)
		if user.Username == userFromDB.Username && user.Password == userFromDB.Password {
			exist = 0

			{
				//return data
			}
			break
		}
	}
	if exist == 1 {
		{
			//return response
		}
	}
	db.Close()
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
func auth(c *gin.Context) {
	var username string
	c.Bind(&username)
	var claims *_userClaims
	claims.username = string(username)
	sign := GenerateToken(claims) //签发token
	c.Set("Authorization", sign)
}
func GenerateToken(claims *_userClaims) string { //生成token
	claims.ExpiresAt = time.Now().Add(EffectTime).Unix()                                //过期时间
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret) //签发token
	CheckErr(err)
	if err != nil {
		return ""
	}
	return sign
}
func JwtVerify(c *gin.Context) { //验证token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Set("Authorization", "none")
	}
	jwtUsername, err := ParseToken(token)
	if err != nil {
		c.Set("Authorization", "err")
	} else {
		c.Set("Authorization", jwtUsername)
	}
}
func ParseToken(tokenGot string) (*_userClaims, error) {
	// 解析token（传参分别为：字符串token，将解析结果保存至指定的结构体，返回生成token时所使用的secret从而用于解签名的方法）
	token, err := jwt.ParseWithClaims(tokenGot, &_userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil { //解析出错 即鉴权失败
		return nil, err
	}
	if token != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		if claims, ok := token.Claims.(*_userClaims); ok && token.Valid {
			return claims, nil
		}
	}
	CheckErr(err)
	return nil, err
}
