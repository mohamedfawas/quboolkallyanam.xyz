package validation

const (
	ConstantHomeDistrictThiruvananthapuram = "thiruvananthapuram"
	ConstantHomeDistrictKollam             = "kollam"
	ConstantHomeDistrictPathanamthitta     = "pathanamthitta"
	ConstantHomeDistrictAlappuzha          = "alappuzha"
	ConstantHomeDistrictKottayam           = "kottayam"
	ConstantHomeDistrictErnakulam          = "ernakulam"
	ConstantHomeDistrictThrissur           = "thrissur"
	ConstantHomeDistrictPalakkad           = "palakkad"
	ConstantHomeDistrictMalappuram         = "malappuram"
	ConstantHomeDistrictKozhikode          = "kozhikode"
	ConstantHomeDistrictWayanad            = "wayanad"
	ConstantHomeDistrictKannur             = "kannur"
	ConstantHomeDistrictKasaragod          = "kasaragod"
	ConstantHomeDistrictIdukki             = "idukki"
	ConstantHomeDistrictNotMentioned       = "not_mentioned"
)

func IsValidHomeDistrict(homeDistrict string) bool {
	switch homeDistrict {
	case ConstantHomeDistrictThiruvananthapuram, ConstantHomeDistrictKollam,
		ConstantHomeDistrictPathanamthitta, ConstantHomeDistrictAlappuzha,
		ConstantHomeDistrictKottayam, ConstantHomeDistrictErnakulam,
		ConstantHomeDistrictThrissur, ConstantHomeDistrictPalakkad,
		ConstantHomeDistrictMalappuram, ConstantHomeDistrictKozhikode,
		ConstantHomeDistrictWayanad, ConstantHomeDistrictKannur,
		ConstantHomeDistrictKasaragod, ConstantHomeDistrictIdukki,
		ConstantHomeDistrictNotMentioned:
		return true
	default:
		return false
	}
}
