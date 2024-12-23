package shelf

import (
	author2 "booktastic-server-go/author"
	"booktastic-server-go/book"
	"booktastic-server-go/database"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type Shelf struct {
	ID           uint64          `json:"id"`
	Externaluid  string          `json:"externaluid"`
	Ouruid       string          `json:"ouruid" gorm:"-"`
	Externalmods json.RawMessage `json:"externalmods"`
	Processed    bool            `json:"processed"`
	Rating       string          `json:"rating"`
	Created      string          `json:"created" gorm:"<-:false"`
}

func Create(c *fiber.Ctx) error {
	var shelf Shelf

	err := json.Unmarshal(c.Body(), &shelf)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	db := database.GetDB()
	db.Create(&shelf)

	if shelf.ID > 0 {
		return c.JSON(shelf)
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Failed to create")
}

func List(c *fiber.Ctx) error {
	var shelves []Shelf

	db := database.GetDB()
	db.Find(&shelves)

	return c.JSON(shelves)
}

func Single(c *fiber.Ctx) error {
	var shelf Shelf
	var found bool

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err == nil {
		db := database.GetDB()

		err := db.Where("id = ?", id).First(&shelf).Error
		found = !errors.Is(err, gorm.ErrRecordNotFound)

		if found {
			return c.JSON(shelf)
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "Not found")
}

func Books(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var books []book.Book

	if err == nil {
		db := database.GetDB()

		db.Raw("SELECT books.id, title, isbn13 FROM `books` "+
			"INNER JOIN shelves_books ON shelves_books.bookid = books.id "+
			"WHERE shelfid = ?;", id).Scan(&books)

		for i := range books {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				var author author2.Author
				db.Raw("SELECT authors.* FROM `authors` INNER JOIN books_authors ON books_authors.authorid = authors.id WHERE bookid = ?;", books[i].ID).Scan(&author)

				mu.Lock()
				books[i].Authors = append(books[i].Authors, author)
				mu.Unlock()
			}(i)
		}

		wg.Wait()
	}

	return c.JSON(books)
}

func Patch(c *fiber.Ctx) error {
	var shelf Shelf
	var found bool

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err == nil {
		db := database.GetDB()

		err := db.Where("id = ?", id).First(&shelf).Error
		found = !errors.Is(err, gorm.ErrRecordNotFound)

		if found {
			// We support updating the processed flag and the rating.
			err = json.Unmarshal(c.Body(), &shelf)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
			}

			if !shelf.Processed {
				db.Exec("DELETE FROM shelves_books WHERE shelfid = ?", id)
			}

			db.Exec("UPDATE shelves SET processed = ? WHERE id = ?", shelf.Processed, id)

			if shelf.Rating != "" {
				db.Exec("UPDATE shelves SET rating = ? WHERE id = ?", shelf.Rating, id)
			}

			return c.JSON(shelf)
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "Not found")
}
