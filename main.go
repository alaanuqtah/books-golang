package main

import (
	"go-fiber-postgres/models"
	"go-fiber-postgres/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type Book struct {
	Author    string `json:author`
	Title     string `json:title`
	Publisher string `json:publisher`
}

type Repo struct {
	DB *gorm.DB
}

func (r *Repo) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book) //fiber is so abstarct where we dont have access to (w or r like in http package) context only lets you see the body thats why we need to parse and decode it from json to struct

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldnt create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book has been added"})
	return nil

}

func (r *Repo) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "books fetched successfuly",
			"data": bookModels})
	return nil

}

func (r *Repo) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Book{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Delete(bookModel, id) //query

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "couldnt delete book",
		})
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "deleted successfuly",
	})
	return nil

}
func (r *Repo) GetBookByID(context *fiber.Ctx) error {
	bookModel := &models.Book{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get book at that id"})
		return nil
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book fetched successfuly",
			"data": bookModel})
	return nil

}

func (r *Repo) SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	//env setup
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config) //create connection with postrges
	if err != nil {
		log.Fatal(err)
	}

	err = models.MigrateBooks(db) //create db

	if err != nil {
		log.Fatal(err)
	}
	r := Repo{
		DB: db,
	}

	//app setup
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8000")

}
