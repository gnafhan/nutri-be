package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BahanMakananController struct {
	BahanMakananService service.BahanMakananService
}

func NewBahanMakananController(service service.BahanMakananService) *BahanMakananController {
	return &BahanMakananController{
		BahanMakananService: service,
	}
}

// @Tags         BahanMakanan
// @Summary      Get all bahan makanan
// @Description  Get all bahan makanan
// @Security     BearerAuth
// @Produce      json
// @Router       /bahan-makanan [get]
// @Success      200  {object}  response.SuccessWithBahanMakananList
func (c *BahanMakananController) GetAllBahanMakanan(ctx *fiber.Ctx) error {
	bahanMakananList, err := c.BahanMakananService.GetAllBahanMakanan(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithBahanMakananList{
		Status:  "success",
		Message: "Bahan makanan fetched successfully",
		Data:    bahanMakananList,
	})
}

// @Tags         BahanMakanan
// @Summary      Get bahan makanan by kode
// @Description  Get bahan makanan by kode
// @Security     BearerAuth
// @Produce      json
// @Param        kode  path  string  true  "Kode Bahan Makanan"
// @Router       /bahan-makanan/kode/{kode} [get]
// @Success      200  {object}  response.SuccessWithBahanMakanan
func (c *BahanMakananController) GetBahanMakananByKode(ctx *fiber.Ctx) error {
	kode := ctx.Params("kode")
	if kode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid kode parameter",
		})
	}

	bahanMakanan, err := c.BahanMakananService.GetBahanMakananByKode(ctx, kode)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithBahanMakanan{
		Status:  "success",
		Message: "Bahan makanan fetched successfully",
		Data:    *bahanMakanan,
	})
}

// @Tags         BahanMakanan
// @Summary      Get bahan makanan by ID
// @Description  Get bahan makanan by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  integer  true  "Bahan Makanan ID"
// @Router       /bahan-makanan/{id} [get]
// @Success      200  {object}  response.SuccessWithBahanMakanan
func (c *BahanMakananController) GetBahanMakananById(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid ID parameter",
		})
	}

	bahanMakanan, err := c.BahanMakananService.GetBahanMakananById(ctx, uint32(id))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithBahanMakanan{
		Status:  "success",
		Message: "Bahan makanan fetched successfully",
		Data:    *bahanMakanan,
	})
}

// @Tags         BahanMakanan
// @Summary      Get bahan makanan by mentah olahan
// @Description  Get bahan makanan by mentah olahan status
// @Security     BearerAuth
// @Produce      json
// @Param        mentah_olahan  path  string  true  "Mentah/Olahan Status"
// @Router       /bahan-makanan/mentah-olahan/{mentah_olahan} [get]
// @Success      200  {object}  response.SuccessWithBahanMakananList
func (c *BahanMakananController) GetBahanMakananByMentahOlahan(ctx *fiber.Ctx) error {
	mentahOlahan := ctx.Params("mentah_olahan")
	if mentahOlahan == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid mentah_olahan parameter",
		})
	}

	bahanMakananList, err := c.BahanMakananService.GetBahanMakananByMentahOlahan(ctx, mentahOlahan)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithBahanMakananList{
		Status:  "success",
		Message: "Bahan makanan fetched successfully",
		Data:    bahanMakananList,
	})
}

// @Tags         BahanMakanan
// @Summary      Get bahan makanan by kelompok
// @Description  Get bahan makanan by kelompok makanan
// @Security     BearerAuth
// @Produce      json
// @Param        kelompok  path  string  true  "Kelompok Makanan"
// @Router       /bahan-makanan/kelompok/{kelompok} [get]
// @Success      200  {object}  response.SuccessWithBahanMakananList
func (c *BahanMakananController) GetBahanMakananByKelompok(ctx *fiber.Ctx) error {
	kelompokMakanan := ctx.Params("kelompok")
	if kelompokMakanan == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid kelompok parameter",
		})
	}

	bahanMakananList, err := c.BahanMakananService.GetBahanMakananByKelompok(ctx, kelompokMakanan)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithBahanMakananList{
		Status:  "success",
		Message: "Bahan makanan fetched successfully",
		Data:    bahanMakananList,
	})
}

// @Tags         BahanMakanan
// @Summary      Update bahan makanan
// @Description  Update bahan makanan
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  integer  true  "Bahan Makanan ID"
// @Param        request  body  model.BahanMakanan  true  "Bahan Makanan data"
// @Router       /bahan-makanan/{id} [put]
// @Success      200  {object}  response.SuccessWithBahanMakanan
func (c *BahanMakananController) UpdateBahanMakanan(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid ID parameter",
		})
	}

	var request model.BahanMakanan
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	bahanMakanan, err := c.BahanMakananService.UpdateBahanMakanan(ctx, uint32(id), &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithBahanMakanan{
		Status:  "success",
		Message: "Bahan makanan updated successfully",
		Data:    *bahanMakanan,
	})
}
