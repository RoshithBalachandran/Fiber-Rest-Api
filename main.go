package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/roshith/fiber-gorm/models"
	"github.com/roshith/fiber-gorm/storage"
	"gorm.io/gorm"
)

type Book struct {
	Auther    string `json:"auther"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create a book"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book hasbeen created"})
	return err
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	booksModel := &[]models.Book{}
	err := r.DB.Find(booksModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not find the book",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book fetched sucessfull",
		"Data":    booksModel,
	})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookmodel := models.Book{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Id cannot be  empty",
		})
		return nil
	}
	err := r.DB.Delete(bookmodel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delet book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book deleted sucessfully",
	})
	return nil
}
func (r *Repository) GetBookId(context *fiber.Ctx) error {
	id := context.Params("id")
	bookmodels := &models.Book{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	fmt.Println("the id is ", id)
	err := r.DB.Where("id=?", id).First(bookmodels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get the book",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book with id fetched sucessfully",
		"Data":    bookmodels,
	})
	return nil

}
func (r *Repository) SetupRouter(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/Create_book", r.CreateBook)
	api.Delete("/delete_book", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookId)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not onnect the databse")

	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db")
	}
	r := Repository{
		DB: db,
	}
	app := fiber.New()

	r.SetupRouter(app)

	app.Listen(":8050")
}
