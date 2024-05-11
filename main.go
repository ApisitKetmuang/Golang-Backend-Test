package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "dev.opensource-technology.com"  
	port     = 5523         
	user     = "posts "     
	password = "38S2GPNZut4Tmvan" 
	dbname   = "posts " 
)

type Post struct {
	ID       	string 		`json:"id"`
	Title 		string 		`json:"title"`
	Content 	string 		`json:"content"`
	Published 	bool 		`json:"published"`
	CreatedAt   time.Time 	`json:"created_at"`
}

var db *sql.DB

func main() {

	psqlInfo := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
	host, 	port, 	user, 	password, 	dbname)

	sdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	db = sdb

	defer db.Close()
	
	err= db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/post", getPostsHandler)
	app.Get("/draft", getDraftsHandler)
	app.Get("/post/:id", getPostHandler)
	// app.Get("/post/:id?/:title?", getPostHandler)

	app.Post("/post", createPostHandler)
	app.Put("/post/:id", updatePostHandler)
	app.Delete("/post/:id", deletePostHandler)


	app.Listen(":8080")
}

//draft
func getDraftsHandler(c *fiber.Ctx) error {
	post, err := getDraft()
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(post)
}

//published
func getPostsHandler(c *fiber.Ctx) error {
	post, err := getPosts()
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(post)
}

//getPostById
func getPostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	// title := c.Params("title")
	post, err := getPost(postId)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(post)
}

//CratePostById
func createPostHandler(c *fiber.Ctx) error {
	p := new(Post)

	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := createPost(p)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	
	return c.JSON(p)
}

//UpdatePostById
func updatePostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	p := new(Post)

	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	post := updatePost(postId, p)
	return c.JSON(post)
}

//DeletePostById
func deletePostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	err := deletePost(postId)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

