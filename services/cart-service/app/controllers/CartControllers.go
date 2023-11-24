package controllers

import (
	models "cart_service/app/models"
	database "cart_service/db"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type CartController struct {
	// Add        func(*fiber.Ctx) error
	// GetById    func(*fiber.Ctx) error
	// GetByOwner func(*fiber.Ctx) error
	// EditById   func(*fiber.Ctx) error
	// DeleteById func(*fiber.Ctx) error
}

func (ctrl CartController) Add(c *fiber.Ctx) error {

	request := models.AddCartRequest{}

	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(403)
	}

	newCart := models.Cart{
		ID:        guuid.New(),
		UserId:    request.UserId,
		ProductId: request.ProductId,
	}

	// add product id to user table for easy searching (?)

	db := database.DB
	db.Create(&newCart)
	return c.SendStatus(200)
}

func (ctrl CartController) GetByUser(c *fiber.Ctx) error {

	userId := c.Params("id")
	db := database.DB
	foundCart := []models.Cart{}
	err := db.Where("user_id = ?", userId).Find(&foundCart).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
		})
	}

	if err != nil {

		return c.SendStatus(401)
	}

	var totalPrice float32 = 0
	var items []models.Product

	for _, item := range foundCart {
		foundProduct := models.Product{}
		err := db.First(&foundProduct, "id = ?", item.ProductId).Error
		if err == nil {
			totalPrice += foundProduct.Price
			items = append(items, foundProduct)
		}
	}

	var cartResponse models.CartResponse
	cartResponse.TotalPrice = totalPrice
	cartResponse.Items = items

	response, responseErr := json.Marshal(cartResponse)
	if responseErr != nil {
		panic(err)
	}

	return c.SendString(string(response))
}

func (ctrl CartController) DeleteById(c *fiber.Ctx) error {

	cartId := c.Params("id")
	db := database.DB
	foundCart := models.Cart{}
	err := db.First(&foundCart, "id = ?", cartId).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
		})
	}

	if err != nil {

		return c.SendStatus(401)
	}

	db.Delete(&foundCart)
	return c.SendStatus(200)
}

func InitializeCartController() CartController {
	return CartController{}
}
