package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
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
type _claims struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}                                              //jwt struct
const TokenExpireDuration = time.Hour * 24 * 2 //过期时间2d
var secret []byte = []byte("SYMCSigned")       //盐
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
	r.Use(JWTMiddleware) //token中间件
	r.Group("/api/post") //推文路由组
	{
		r.POST("/", func(context *gin.Context) {

		})

		r.POST("/addPost", AddPost)
		r.GET("/getPost", GetPost)
		r.PUT("/putPost", PutPost)
		r.DELETE("/api/post/deletePost", DeletePost)
	}

	r.Group("/api/user") //用户路由组
	{
		r.POST("/register", Register)
		r.POST("/login", Login)
		r.POST("/auth/token", GetTokenHandler)
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
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Create(&post)
	postNums++
	{
		//return response
	}
}
func GetPost(c *gin.Context) { //获取推文
	//随机一条推文
	var post _post
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Take(&post)
	{
		//return data
	}
}
func PutPost(c *gin.Context) { //更新推文
	var post _post
	c.ShouldBindJSON(&post)
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Where("id=?", post.Id).Update("content", post.Content)
	{
		//return response
	}
}
func DeletePost(c *gin.Context) { //删除推文
	id := c.Query("Id")
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Where("id=?", id).Delete(&_post{})
	postNums--
	{
		//return response
	}
}

func Register(c *gin.Context) { //注册
	var user _user
	c.ShouldBindJSON(&user)
	var userFromDB _user
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Where("username=?", user.Username).First(&userFromDB)
	//检测名称重复
	if userFromDB.Username != "" {
		//return failed
	} else {
		db.Create(&user)
	}
}
func Login(c *gin.Context) { //登录
	var user _user
	c.ShouldBindJSON(&user)
	var userFromDB _user
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Where("username=?", user.Username).First(&userFromDB)
	if userFromDB.Username != "" {
		//return username
	} else {
		//return failed
	}
}

func AddComment(c *gin.Context) { //添加评论
	var comment _comment
	c.ShouldBindJSON(&comment)
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Create(&comment)
	{
		//return response
	}
}
func GetComment(c *gin.Context) { //获取评论
	var comments []_comment
	var pid int64
	c.Bind(pid)
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Where("pid=?", pid).Find(&comments)
	if comments == nil {
		//return response
	} else {
		//return data
	}
}
func DeleteComment(c *gin.Context) { //删除评论
	var id int64
	c.Bind(&id)
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.Where("id=?", id).Delete(&_comment{})
	commentNums--
	{
		//return response
	}
}
func InitDB() {
	db, err := gorm.Open("sqlite3", "wall.db")
	CheckErr(err)
	defer db.Close()
	db.AutoMigrate(_post{})
	db.AutoMigrate(_comment{})
	db.AutoMigrate(_user{})
}
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GenToken(claim _claims) (string, error) {
	claim.ExpiresAt = time.Now().Add(TokenExpireDuration).Unix() //过期时间
	claim.Issuer = "SYMC"                                        //签发人
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)    //加密
	SignedToken, err := token.SignedString(secret)               //加盐
	CheckErr(err)
	return SignedToken, err
}
func ParseToken(tokenStr string) (*_claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &_claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	CheckErr(err)
	if claims, ok := token.Claims.(*_claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
func JWTMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		c.Abort() //停止
		return
	}
	//解析
	tokenP, err := ParseToken(token)
	if err != nil {
		c.Abort()
		return
	}
	c.Set("Uid", tokenP.Uid)
	c.Set("Username", tokenP.Username)
	c.Next()
}
func GetTokenHandler(c *gin.Context) {
	uid, _ := c.Get("Uid")
	username, _ := c.Get("Username")
	if uid != "" {
		return
		//already has token
	}
	var user _user
	var claims _claims
	c.ShouldBindJSON(&user) //接收
	claims.Uid = user.Uid
	claims.Username = user.Username
	tokenString, err := GenToken(claims)
	CheckErr(err)
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "success",
		"token":    tokenString,
		"uid":      uid,
		"username": username,
	})
	return
}
