package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Article struct
type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Initialize fiber
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize Database
	ConnectDB()

	// Initialize Routes
	app.Post("/api/articles", func(c *fiber.Ctx) error {
		// Bind the request body to an article struct
		article := new(Article)

		if err := c.BodyParser(&article); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Save the article to the database
		Db.Create(&article)

		// Respond with the created article
		return c.Status(fiber.StatusCreated).JSON(article)
	})

	app.Get("/api/articles", func(c *fiber.Ctx) error {
		var articles []Article
		db := Db

		// Get articles from database
		db.Find(&articles)

		return c.Status(fiber.StatusOK).JSON(articles)
	})

	app.Get("/api/articles/:id", func(c *fiber.Ctx) error {
		var article Article
		db := Db

		paramId := c.Params("id")
		id, err := strconv.Atoi(paramId)

		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse id",
			})
			return err
		}

		// Get article from database
		db.Find(&article, id)

		// Return result
		if int(article.ID) == id {
			return c.Status(fiber.StatusOK).JSON(article)
		}

		//
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Article not found",
		})
	})

	app.Put("/api/articles/:id", func(c *fiber.Ctx) error {
		db := Db

		type request struct {
			Title   *string `json:"title"`
			Content *string `json:"content"`
		}

		paramId := c.Params("id")
		id, err := strconv.Atoi(paramId)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse id",
			})
		}

		var body request

		err = c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse body",
			})
		}

		var article Article
		db.First(&article, id)

		// Handle 404
		if int(article.ID) != id {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Article not found",
			})
		}

		if body.Title != nil {
			article.Title = *body.Title
		}

		if body.Content != nil {
			article.Content = *body.Content
		}

		db.Save(&article)

		return c.Status(fiber.StatusOK).JSON(article)
	})

	app.Delete("/api/articles/:id", func(c *fiber.Ctx) error {
		db := Db

		paramId := c.Params("id")

		id, err := strconv.Atoi(paramId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse id",
			})
		}

		var article Article
		db.First(&article, id)

		if int(article.ID) != id {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Article not found.",
			})
		}

		db.Delete(&article)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "Artice deleted successfully",
		})
	})

	// Start Server
	app.Listen(":3000")
}

var Db *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database.")
	}

	// Migrate the schema
	db.AutoMigrate(&Article{})

	Db = db
}
