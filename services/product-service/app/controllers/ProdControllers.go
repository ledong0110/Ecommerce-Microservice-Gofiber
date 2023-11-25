package controllers

import (
	"encoding/json"
	models "product_service/app/models"
	database "product_service/db"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type ProdController struct {
	// Add        func(*fiber.Ctx) error
	// GetById    func(*fiber.Ctx) error
	// GetByOwner func(*fiber.Ctx) error
	// EditById   func(*fiber.Ctx) error
	// DeleteById func(*fiber.Ctx) error
}

// Add godoc
//
//	@Summary		Post new product
//	@Description	Post new product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body models.ProductForm true "New Product Form"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/prod/Add/ [post]
func (ctrl ProdController) Add(c *fiber.Ctx) error {

	request := models.ProductForm{}

	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(403)
	}

	newProduct := models.Product{
		ID:      guuid.New(),
		Name:    request.Name,
		Price:   request.Price,
		Picture: request.Picture,
		Detail:  request.Detail,
		OwnerId: request.UserId,
	}

	// add product id to user table for easy searching (?)

	db := database.DB
	db.Create(&newProduct)
	return c.SendStatus(200)
}

// Add godoc
//
//	@Summary		Get product by id
//	@Description	Get product by id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Product Id"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/prod/GetById/{id} [get]
func (ctrl ProdController) GetById(c *fiber.Ctx) error {

	productId := c.Params("id")
	db := database.DB
	foundProduct := models.Product{}
	err := db.First(&foundProduct, "id = ?", productId).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
		})
	}

	if err != nil {

		return c.SendStatus(401)
	}

	response, responseErr := json.Marshal(foundProduct)
	if responseErr != nil {
		panic(err)
	}

	return c.SendString(string(response))
}

// Add godoc
//
//	@Summary		Get product by owner id
//	@Description	Get product by owner id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Owner Id"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/prod/GetByOwner/{id} [get]
func (ctrl ProdController) GetByOwner(c *fiber.Ctx) error {

	ownerId := c.Params("id")
	db := database.DB
	foundProduct := []models.Product{}
	err := db.Where("owner_id = ?", ownerId).Find(&foundProduct).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
		})
	}

	if err != nil {

		return c.SendStatus(401)
	}

	response, responseErr := json.Marshal(foundProduct)
	if responseErr != nil {
		panic(err)
	}

	return c.SendString(string(response))
}

// Add godoc
//
//	@Summary		Edit existing product by id
//	@Description	Edit existing product by id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Product Id"
//	@Param			request	body models.ProductForm true "Edit Product Form"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/prod/EditById/{id} [put]
func (ctrl ProdController) EditById(c *fiber.Ctx) error {

	request := models.ProductForm{}
	productId := c.Params("id")

	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(403)
	}

	db := database.DB
	foundProduct := models.Product{}
	err := db.First(&foundProduct, "id = ?", productId).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
		})
	}

	if err != nil {

		return c.SendStatus(401)
	}

	foundProduct.Name = request.Name
	foundProduct.Price = request.Price
	foundProduct.Picture = request.Picture
	foundProduct.Detail = request.Detail
	db.Save(&foundProduct)

	return c.SendStatus(200)
}

// Add godoc
//
//	@Summary		Delete product by id
//	@Description	Delete product by id
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Product Id"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/prod/DeleteById/{id} [delete]
func (ctrl ProdController) DeleteById(c *fiber.Ctx) error {

	productId := c.Params("id")
	db := database.DB
	foundProduct := models.Product{}
	err := db.First(&foundProduct, "id = ?", productId).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
		})
	}

	if err != nil {

		return c.SendStatus(401)
	}

	db.Delete(&foundProduct)
	return c.SendStatus(200)
}

func InitializeProdController() ProdController {
	return ProdController{}
}
