package model

type BahanMakanan struct {
	ID                uint32   `json:"id" gorm:"primaryKey"`
	Kode              string   `json:"kode" gorm:"unique"`
	NamaBahanMakanan  string   `json:"nama_bahan_makanan"`
	AirG              float64  `json:"air_g"`
	EnergiKal         float64  `json:"energi_kal"`
	ProteinG          float64  `json:"protein_g"`
	LemakG            float64  `json:"lemak_g"`
	KarbohidratG      float64  `json:"karbohidrat_g"`
	SeratG            *float64 `json:"serat_g"`
	AbuG              float64  `json:"abu_g"`
	KalsiumCaMg       *float64 `json:"kalsium_ca_mg"`
	FosforPMg         *float64 `json:"fosfor_p_mg"`
	BesiFeMg          *float64 `json:"besi_fe_mg"`
	NatriumNaMg       *float64 `json:"natrium_na_mg"`
	KaliumKaMg        *float64 `json:"kalium_ka_mg"`
	TembagaCuMg       *float64 `json:"tembaga_cu_mg"`
	SengZnMg          *float64 `json:"seng_zn_mg"`
	RetinolVitAMcg    *float64 `json:"retinol_vit_a_mcg"`
	BetaKarotenMcg    *float64 `json:"beta_karoten_mcg"`
	KarotenTotalMcg   *float64 `json:"karoten_total_mcg"`
	ThiaminVitB1Mg    *float64 `json:"thiamin_vit_b1_mg"`
	RiboflavinVitB2Mg *float64 `json:"riboflavin_vit_b2_mg"`
	NiasinMg          *float64 `json:"niasin_mg"`
	VitaminCMg        *float64 `json:"vitamin_c_mg"`
	BddPersen         float64  `json:"bdd_persen"`
	MentahOlahan      string   `json:"mentah_olahan"`
	KelompokMakanan   string   `json:"kelompok_makanan"`
}
