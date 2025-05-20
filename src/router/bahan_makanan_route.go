package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func BahanMakananRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, bahanMakananService service.BahanMakananService) {
	bahanMakananController := controller.NewBahanMakananController(bahanMakananService)

	bahanMakanan := v1.Group("/bahan-makanan")
	bahanMakanan.Get("/", m.Auth(u, p), bahanMakananController.GetAllBahanMakanan)
	bahanMakanan.Get("/:id", m.Auth(u, p), bahanMakananController.GetBahanMakananById)
	bahanMakanan.Get("/kode/:kode", m.Auth(u, p), bahanMakananController.GetBahanMakananByKode)
	bahanMakanan.Get("/mentah-olahan/:mentah_olahan", m.Auth(u, p), bahanMakananController.GetBahanMakananByMentahOlahan)
	bahanMakanan.Get("/kelompok/:kelompok", m.Auth(u, p), bahanMakananController.GetBahanMakananByKelompok)
	bahanMakanan.Put("/:id", m.Auth(u, p, "manageUsers"), bahanMakananController.UpdateBahanMakanan)
}
