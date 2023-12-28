// main.go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

// Item model
type Item struct {
	gorm.Model
	Name         string  `json:"name"`
	UnitPrice    float64 `json:"unit_price"`
	ItemCategory string  `json:"item_category"`
}

// Invoice model
type Invoice struct {
	gorm.Model
	Name        string `json:"name"`
	MobileNo    string `json:"mobile_no"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	BillingType string `json:"billing_type"`
}

func main() {
	// Connect to the PostgreSQL database
	db, err = gorm.Open("postgres", "host=localhost user=abcpharmacyuser dbname=abcpharmacy sslmode=disable password=your_password")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Automigrate the database tables
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Invoice{})

	// Initialize the Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to ABC Pharmacy"})
	})

	// Routes for managing items
	itemController := ItemController{}
	app.Get("/items", itemController.GetItems)
	app.Post("/items", itemController.CreateItem)
	app.Put("/items/:id", itemController.UpdateItem)
	app.Delete("/items/:id", itemController.DeleteItem)

	// Route for creating invoices
	invoiceController := InvoiceController{}
	app.Post("/invoices", invoiceController.CreateInvoice)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}

// Controller for managing items
type ItemController struct{}

func (c *ItemController) GetItems(ctx *fiber.Ctx) error {
	var items []Item
	db.Find(&items)
	return ctx.JSON(items)
}

func (c *ItemController) CreateItem(ctx *fiber.Ctx) error {
	var newItem Item
	if err := ctx.BodyParser(&newItem); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&newItem)
	return ctx.Status(fiber.StatusCreated).JSON(newItem)
}

func (c *ItemController) UpdateItem(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var item Item
	if err := db.First(&item, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	var updatedItem Item
	if err := ctx.BodyParser(&updatedItem); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	item.Name = updatedItem.Name
	item.UnitPrice = updatedItem.UnitPrice
	item.ItemCategory = updatedItem.ItemCategory

	db.Save(&item)
	return ctx.Status(fiber.StatusOK).JSON(item)
}

func (c *ItemController) DeleteItem(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var item Item
	if err := db.First(&item, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	db.Delete(&item)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item deleted successfully"})
}

// Controller for creating invoices
type InvoiceController struct{}

func (c *InvoiceController) CreateInvoice(ctx *fiber.Ctx) error {
	var newInvoice Invoice
	if err := ctx.BodyParser(&newInvoice); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&newInvoice)
	return ctx.Status(fiber.StatusCreated).JSON(newInvoice)
}
