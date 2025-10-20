package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func BahanMakananRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, ss service.SubscriptionService, bahanMakananService service.BahanMakananService) {
	bahanMakananController := controller.NewBahanMakananController(bahanMakananService)

	bahanMakanan := v1.Group("/bahan-makanan")
	bahanMakanan.Get("/", m.FreemiumOrAccess(u, p, ss), bahanMakananController.GetAllBahanMakanan)
	bahanMakanan.Get("/:id", m.FreemiumOrAccess(u, p, ss), bahanMakananController.GetBahanMakananById)
	bahanMakanan.Get("/kode/:kode", m.FreemiumOrAccess(u, p, ss), bahanMakananController.GetBahanMakananByKode)
	bahanMakanan.Get("/mentah-olahan/:mentah_olahan", m.FreemiumOrAccess(u, p, ss), bahanMakananController.GetBahanMakananByMentahOlahan)
	bahanMakanan.Get("/kelompok/:kelompok", m.FreemiumOrAccess(u, p, ss), bahanMakananController.GetBahanMakananByKelompok)
	bahanMakanan.Put("/:id", m.Auth(u, p, "manageUsers"), bahanMakananController.UpdateBahanMakanan)
}
