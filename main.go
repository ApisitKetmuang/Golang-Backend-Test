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
	host, port, user, password, dbname)

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

	app.Get("/posts", getPostsHandler)
	app.Get("/posts/draft", getDraftsHandler)
	app.Get("/posts/:id", getPostHandler)

	app.Post("/posts", createPostHandler)
	app.Put("/posts/:id", updatePostHandler)
	app.Patch("/posts/:id", publishedPostHandler)
	app.Delete("/posts/:id", deletePostHandler)

	app.Listen(":8080")
}

//DraftPage
func getDraftsHandler(c *fiber.Ctx) error {
	post, err := getDraft()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(200).JSON(fiber.Map{"posts": post})
}

//PublishedPage
func getPostsHandler(c *fiber.Ctx) error {
	post, err := getPosts()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{"posts": post})
}

//GetPostById
func getPostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	post, err := getPost(postId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(post)
}

//CreatePost
func createPostHandler(c *fiber.Ctx) error {
	post := new(Post)
	
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	
	if post.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title is required"})
	}
	
	err := createPost(post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(201).JSON(fiber.Map{"message": "Post has created", "post": post})
}

//UpdatePostById
func updatePostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	p := new(Post)

	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	post := updatePost(postId, p)
	return c.JSON(post)
}

func publishedPostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	p := new(Post)

	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	post := publishedPost(postId, p)
	return c.Status(201).JSON(fiber.Map{"message": "Post has published", "post": post})
}

//DeletePostById
func deletePostHandler(c *fiber.Ctx) error {
	postId := c.Params("id")
	err := deletePost(postId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(204).JSON(fiber.Map{"message": "Post deleted"})
}

