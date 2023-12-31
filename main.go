package main

import (
	"fmt"
	"os"
	"time"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID          string `json:"id" gorm:"primaryKey"`
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoURL"`
	Class       string `json:"class"`
	Faculty     string `json:"faculty"`
	Department  string `json:"department"`
	Grade       string `json:"grade"`
	Can         string `json:"can"`
	Did         string `json:"did"`
	Will        string `json:"will"`
	IsPublic    bool   `json:"isPublic"`
}

type Author struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	AvatarURL string
}

type Comment struct {
	ID     uint
	PostID uint
	Name   string
	Message string
	PostedAt time.Time
}

type Post struct {
	gorm.Model
	Title           string           `json:"title"`
	AuthorID        string           `json:"authorId"`
	// Tags            string           `gorm:"type:text" json:"tags"`
	Category        string           `json:"category"`
	Tech 		    string           `json:"tech"`
	Curriculum 	    string           `json:"curriculum"`
	Content         string           `json:"content"`
	CoverURL        string           `json:"coverUrl"`
	MetaTitle       string           `json:"metaTitle"`
	TotalViews      int              `json:"totalViews"`
	TotalShares     int              `json:"totalShares"`
	Description     string           `json:"description"`
	TotalComments   int              `json:"totalComments"`
	TotalFavorites  int              `json:"totalFavorites"`
	// MetaKeywords    string           `gorm:"type:text" json:"metaKeywords"`
	// Comments        []Comment        `gorm:"foreignKey:PostID" json:"comments"`  // No change, as comments field is commented out in your IPostItem type
	Author          Author           `gorm:"foreignKey:ID" json:"author"`
}


func main() {
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
	dbSocketDir := os.Getenv("DB_HOST")
    // instanceConnectionName := os.Getenv("DB_HOST")  
    dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true", dbUser, dbPass, dbSocketDir, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&Post{}, &Author{}, &User{})

	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	

	//READ CRUD create read update delete

	r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello, World!")
    })

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},  // 全てのオリジンからのアクセスを許可
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept","Authorization"},
        AllowCredentials: true,
    }))
	r.POST("/create-user", func(c *gin.Context) {
		var newUser User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		  }
		
		  // GORMを使って新しいユーザーをデータベースに保存
		  result := db.Create(&newUser)
		  if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		  }
		
		  c.JSON(200, gin.H{"data": newUser})
		})
		

	r.PUT("/update-user/:id", func(c *gin.Context) {
		var user User
		userID := c.Param("id")
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		  }
		
		  // IDで既存のユーザーを検索し、情報を更新
		  result := db.Model(&User{}).Where("id = ?", userID).Updates(user)
		  if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		  }
		
		  c.JSON(200, gin.H{"data": user})
		})

	r.GET("/user/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		log.Printf("Requested id: %s\n", id)
		err := db.Where("id = ?", id).First(&user).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user": user})
	})

	r.DELETE("/delete-user/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id") // URLからタイトルを取得
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			// 該当する投稿が見つからない場合、404エラーを返す
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}
	
		// GORMのDeleteメソッドを使用して投稿を削除
		if err := db.Delete(&user).Error; err != nil {
			// データベースエラーが発生した場合、500エラーを返す
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	
		// 削除に成功したら、200ステータスとメッセージを返す
		c.JSON(200, gin.H{"message": "User deleted successfully!"})
	})
		

	r.GET("/posts", func(c *gin.Context) {
		var posts []Post
		err := db.Find(&posts).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"posts": posts})
	})

	r.GET("/posts/:id", func(c *gin.Context) {
		var post Post
		id := c.Param("id")
		log.Printf("Requested id: %s\n", id)
		err := db.Where("id = ?", id).First(&post).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"post": post})
	})
	
  
	r.POST("/create-posts", func(c *gin.Context) {
		var newPost []Post
		if err := c.ShouldBindJSON(&newPost); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := db.Create(&newPost).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Posts created successfully!", "posts": newPost})
	})

	r.POST("/create-post", func(c *gin.Context) {
		var newPost Post
		if err := c.ShouldBindJSON(&newPost); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := db.Create(&newPost).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Post created successfully!", "post": newPost})
	})

	r.GET("/search", func(c *gin.Context) {
		var posts []Post
		query := c.Query("query") // クエリパラメータ 'query' を取得

		// タイトル、内容、または説明にクエリが含まれる投稿を検索
		result := db.Where("title LIKE ? OR content LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&posts)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(200, gin.H{"results": posts})
	})

	r.PUT("/edit/:id", func(c *gin.Context) {
		id := c.Param("id")
		var post Post
		if err := db.Where("id = ?", id).First(&post).Error; err != nil {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}

		// リクエストボディから更新データをバインド
		var updateData Post
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// GORMのUpdatesを使って、バインドされたデータで既存のレコードを更新
		if err := db.Model(&post).Updates(updateData).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Post updated successfully!", "post": post})
	})

	r.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id") // URLからタイトルを取得
		var post Post
		if err := db.Where("id = ?", id).First(&post).Error; err != nil {
			// 該当する投稿が見つからない場合、404エラーを返す
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}
	
		// GORMのDeleteメソッドを使用して投稿を削除
		if err := db.Delete(&post).Error; err != nil {
			// データベースエラーが発生した場合、500エラーを返す
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	
		// 削除に成功したら、200ステータスとメッセージを返す
		c.JSON(200, gin.H{"message": "Post deleted successfully!"})
	})
	
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // デフォルトポート
    }
    r.Run(":" + port)

    log.Println("Server started on port " + os.Getenv("PORT"))
    // r.Run(":8080")

}


