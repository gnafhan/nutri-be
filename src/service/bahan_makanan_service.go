package service

import (
	"app/src/grpc"
	pb "app/src/grpc/proto/bahan_makanan"
	"app/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type BahanMakananService interface {
	GetAllBahanMakanan(ctx *fiber.Ctx) ([]model.BahanMakanan, error)
	GetBahanMakananByKode(ctx *fiber.Ctx, kode string) (*model.BahanMakanan, error)
	GetBahanMakananById(ctx *fiber.Ctx, id uint32) (*model.BahanMakanan, error)
	GetBahanMakananByMentahOlahan(ctx *fiber.Ctx, mentahOlahan string) ([]model.BahanMakanan, error)
	GetBahanMakananByKelompok(ctx *fiber.Ctx, kelompokMakanan string) ([]model.BahanMakanan, error)
	UpdateBahanMakanan(ctx *fiber.Ctx, id uint32, bahanMakanan *model.BahanMakanan) (*model.BahanMakanan, error)
}

type bahanMakananService struct {
	Log    *logrus.Logger
	Client *grpc.BahanMakananClient
}

func NewBahanMakananService(client *grpc.BahanMakananClient) BahanMakananService {
	return &bahanMakananService{
		Log:    logrus.New(),
		Client: client,
	}
}

func ConvertPbToModel(pb *pb.BahanMakanan) model.BahanMakanan {
	bahanMakanan := model.BahanMakanan{
		ID:               pb.Id,
		Kode:             pb.Kode,
		NamaBahanMakanan: pb.NamaBahanMakanan,
		AirG:             pb.AirG,
		EnergiKal:        pb.EnergiKal,
		ProteinG:         pb.ProteinG,
		LemakG:           pb.LemakG,
		KarbohidratG:     pb.KarbohidratG,
		AbuG:             pb.AbuG,
		BddPersen:        pb.BddPersen,
		MentahOlahan:     pb.MentahOlahan,
		KelompokMakanan:  pb.KelompokMakanan,
	}

	if pb.SeratG != nil {
		bahanMakanan.SeratG = pb.SeratG
	}
	if pb.KalsiumCaMg != nil {
		bahanMakanan.KalsiumCaMg = pb.KalsiumCaMg
	}
	if pb.FosforPMg != nil {
		bahanMakanan.FosforPMg = pb.FosforPMg
	}
	if pb.BesiFeMg != nil {
		bahanMakanan.BesiFeMg = pb.BesiFeMg
	}
	if pb.NatriumNaMg != nil {
		bahanMakanan.NatriumNaMg = pb.NatriumNaMg
	}
	if pb.KaliumKaMg != nil {
		bahanMakanan.KaliumKaMg = pb.KaliumKaMg
	}
	if pb.TembagaCuMg != nil {
		bahanMakanan.TembagaCuMg = pb.TembagaCuMg
	}
	if pb.SengZnMg != nil {
		bahanMakanan.SengZnMg = pb.SengZnMg
	}
	if pb.RetinolVitAMcg != nil {
		bahanMakanan.RetinolVitAMcg = pb.RetinolVitAMcg
	}
	if pb.BetaKarotenMcg != nil {
		bahanMakanan.BetaKarotenMcg = pb.BetaKarotenMcg
	}
	if pb.KarotenTotalMcg != nil {
		bahanMakanan.KarotenTotalMcg = pb.KarotenTotalMcg
	}
	if pb.ThiaminVitB1Mg != nil {
		bahanMakanan.ThiaminVitB1Mg = pb.ThiaminVitB1Mg
	}
	if pb.RiboflavinVitB2Mg != nil {
		bahanMakanan.RiboflavinVitB2Mg = pb.RiboflavinVitB2Mg
	}
	if pb.NiasinMg != nil {
		bahanMakanan.NiasinMg = pb.NiasinMg
	}
	if pb.VitaminCMg != nil {
		bahanMakanan.VitaminCMg = pb.VitaminCMg
	}

	return bahanMakanan
}

func ConvertModelToPb(model *model.BahanMakanan) *pb.BahanMakanan {
	pbBahanMakanan := &pb.BahanMakanan{
		Id:               model.ID,
		Kode:             model.Kode,
		NamaBahanMakanan: model.NamaBahanMakanan,
		AirG:             model.AirG,
		EnergiKal:        model.EnergiKal,
		ProteinG:         model.ProteinG,
		LemakG:           model.LemakG,
		KarbohidratG:     model.KarbohidratG,
		AbuG:             model.AbuG,
		BddPersen:        model.BddPersen,
		MentahOlahan:     model.MentahOlahan,
		KelompokMakanan:  model.KelompokMakanan,
	}

	if model.SeratG != nil {
		pbBahanMakanan.SeratG = *&model.SeratG
	}
	if model.KalsiumCaMg != nil {
		pbBahanMakanan.KalsiumCaMg = *&model.KalsiumCaMg
	}
	if model.FosforPMg != nil {
		pbBahanMakanan.FosforPMg = *&model.FosforPMg
	}
	if model.BesiFeMg != nil {
		pbBahanMakanan.BesiFeMg = *&model.BesiFeMg
	}
	if model.NatriumNaMg != nil {
		pbBahanMakanan.NatriumNaMg = *&model.NatriumNaMg
	}
	if model.KaliumKaMg != nil {
		pbBahanMakanan.KaliumKaMg = model.KaliumKaMg
	}
	if model.TembagaCuMg != nil {
		pbBahanMakanan.TembagaCuMg = model.TembagaCuMg
	}
	if model.SengZnMg != nil {
		pbBahanMakanan.SengZnMg = model.SengZnMg
	}
	if model.RetinolVitAMcg != nil {
		pbBahanMakanan.RetinolVitAMcg = model.RetinolVitAMcg
	}
	if model.BetaKarotenMcg != nil {
		pbBahanMakanan.BetaKarotenMcg = model.BetaKarotenMcg
	}
	if model.KarotenTotalMcg != nil {
		pbBahanMakanan.KarotenTotalMcg = model.KarotenTotalMcg
	}
	if model.ThiaminVitB1Mg != nil {
		pbBahanMakanan.ThiaminVitB1Mg = model.ThiaminVitB1Mg
	}
	if model.RiboflavinVitB2Mg != nil {
		pbBahanMakanan.RiboflavinVitB2Mg = model.RiboflavinVitB2Mg
	}
	if model.NiasinMg != nil {
		pbBahanMakanan.NiasinMg = model.NiasinMg
	}
	if model.VitaminCMg != nil {
		pbBahanMakanan.VitaminCMg = model.VitaminCMg
	}

	return pbBahanMakanan
}

func (s *bahanMakananService) GetAllBahanMakanan(ctx *fiber.Ctx) ([]model.BahanMakanan, error) {
	response, err := s.Client.GetAllBahanMakanan(ctx.Context())
	if err != nil {
		s.Log.Errorf("Failed to get all bahan makanan: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get bahan makanan data")
	}

	var bahanMakananList []model.BahanMakanan
	for _, pbBahanMakanan := range response.BahanMakanan {
		bahanMakananList = append(bahanMakananList, ConvertPbToModel(pbBahanMakanan))
	}

	return bahanMakananList, nil
}

func (s *bahanMakananService) GetBahanMakananByKode(ctx *fiber.Ctx, kode string) (*model.BahanMakanan, error) {
	response, err := s.Client.GetBahanMakananByKode(ctx.Context(), kode)
	if err != nil {
		s.Log.Errorf("Failed to get bahan makanan by kode: %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Bahan makanan not found")
	}

	bahanMakanan := ConvertPbToModel(response.BahanMakanan)
	return &bahanMakanan, nil
}

func (s *bahanMakananService) GetBahanMakananById(ctx *fiber.Ctx, id uint32) (*model.BahanMakanan, error) {
	response, err := s.Client.GetBahanMakananById(ctx.Context(), id)
	if err != nil {
		s.Log.Errorf("Failed to get bahan makanan by id: %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Bahan makanan not found")
	}

	bahanMakanan := ConvertPbToModel(response.BahanMakanan)
	return &bahanMakanan, nil
}

func (s *bahanMakananService) GetBahanMakananByMentahOlahan(ctx *fiber.Ctx, mentahOlahan string) ([]model.BahanMakanan, error) {
	response, err := s.Client.GetBahanMakananByMentahOlahan(ctx.Context(), mentahOlahan)
	if err != nil {
		s.Log.Errorf("Failed to get bahan makanan by mentah/olahan: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get bahan makanan data")
	}

	var bahanMakananList []model.BahanMakanan
	for _, pbBahanMakanan := range response.BahanMakanan {
		bahanMakananList = append(bahanMakananList, ConvertPbToModel(pbBahanMakanan))
	}

	return bahanMakananList, nil
}

func (s *bahanMakananService) GetBahanMakananByKelompok(ctx *fiber.Ctx, kelompokMakanan string) ([]model.BahanMakanan, error) {
	response, err := s.Client.GetBahanMakananByKelompok(ctx.Context(), kelompokMakanan)
	if err != nil {
		s.Log.Errorf("Failed to get bahan makanan by kelompok: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get bahan makanan data")
	}

	var bahanMakananList []model.BahanMakanan
	for _, pbBahanMakanan := range response.BahanMakanan {
		bahanMakananList = append(bahanMakananList, ConvertPbToModel(pbBahanMakanan))
	}

	return bahanMakananList, nil
}

func (s *bahanMakananService) UpdateBahanMakanan(ctx *fiber.Ctx, id uint32, bahanMakanan *model.BahanMakanan) (*model.BahanMakanan, error) {
	pbBahanMakanan := ConvertModelToPb(bahanMakanan)

	response, err := s.Client.UpdateBahanMakanan(ctx.Context(), id, pbBahanMakanan)
	if err != nil {
		s.Log.Errorf("Failed to update bahan makanan: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update bahan makanan")
	}

	updatedBahanMakanan := ConvertPbToModel(response.BahanMakanan)
	return &updatedBahanMakanan, nil
}
