package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func BahanMakananRoutes(v1 fiber.Router, u service.UserService, ss service.SubscriptionService, bahanMakananService service.BahanMakananService) {
	bahanMakananController := controller.NewBahanMakananController(bahanMakananService)

	bahanMakanan := v1.Group("/bahan-makanan")
	bahanMakanan.Get("/", m.FreemiumOrAccess(u, nil, ss), bahanMakananController.GetAllBahanMakanan)
	bahanMakanan.Get("/:id", m.FreemiumOrAccess(u, nil, ss), bahanMakananController.GetBahanMakananById)
	bahanMakanan.Get("/kode/:kode", m.FreemiumOrAccess(u, nil, ss), bahanMakananController.GetBahanMakananByKode)
	bahanMakanan.Get("/mentah-olahan/:mentah_olahan", m.FreemiumOrAccess(u, nil, ss), bahanMakananController.GetBahanMakananByMentahOlahan)
	bahanMakanan.Get("/kelompok/:kelompok", m.FreemiumOrAccess(u, nil, ss), bahanMakananController.GetBahanMakananByKelompok)
	bahanMakanan.Put("/:id", m.Auth(u, nil, "manageUsers"), bahanMakananController.UpdateBahanMakanan)
}
