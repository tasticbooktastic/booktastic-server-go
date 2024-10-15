package shelf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tasticbooktastic/booktastic-server-go/database"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

func (Shelf) TableName() string {
	return "shelf_images"
}

type Shelf struct {
	ID           uint64          `json:"id"`
	Externaluid  string          `json:"externaluid"`
	Ouruid       string          `json:"ouruid"`
	Externalmods json.RawMessage `json:"externalmods"`
	Processed    bool            `json:"processed"`
}

func Single(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	var shelf Shelf
	var found bool

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err == nil {
		db := database.DBConn

		wg.Add(1)

		go func() {
			defer wg.Done()
			err := db.Where("id = ?", id).First(&shelf).Error
			found = !errors.Is(err, gorm.ErrRecordNotFound)
		}()

		wg.Wait()

		if found {
			return c.JSON(shelf)
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "Not found")
}

func Create(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	var shelf Shelf

	err := json.Unmarshal(c.Body(), &shelf)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	db := database.DBConn

	fmt.Printf("Creating shelf %+v\n", shelf)

	wg.Add(1)

	go func() {
		defer wg.Done()
		db.Debug().Create(&shelf)
	}()

	wg.Wait()

	if shelf.ID > 0 {
		return c.JSON(shelf)
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Failed to create")
}
