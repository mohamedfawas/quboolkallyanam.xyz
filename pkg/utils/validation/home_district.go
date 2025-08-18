package validation

import "fmt"

type HomeDistrict string

const (
	Thiruvananthapuram HomeDistrict = "thiruvananthapuram"
	Kollam             HomeDistrict = "kollam"
	Pathanamthitta     HomeDistrict = "pathanamthitta"
	Alappuzha          HomeDistrict = "alappuzha"
	Kottayam           HomeDistrict = "kottayam"
	Ernakulam          HomeDistrict = "ernakulam"
	Thrissur           HomeDistrict = "thrissur"
	Palakkad           HomeDistrict = "palakkad"
	Malappuram         HomeDistrict = "malappuram"
	Kozhikode          HomeDistrict = "kozhikode"
	Wayanad            HomeDistrict = "wayanad"
	Kannur             HomeDistrict = "kannur"
	Kasaragod          HomeDistrict = "kasaragod"
	Idukki             HomeDistrict = "idukki"
	HomeDistrictAny    HomeDistrict = "any"
)

var validHomeDistricts = func() map[HomeDistrict]struct{} {
	vals := []HomeDistrict{
		Thiruvananthapuram, Kollam, Pathanamthitta, Alappuzha,
		Kottayam, Ernakulam, Thrissur, Palakkad,
		Malappuram, Kozhikode, Wayanad, Kannur,
		Kasaragod, Idukki, HomeDistrictAny,
	}
	m := make(map[HomeDistrict]struct{}, len(vals))
	for _, d := range vals {
		m[d] = struct{}{}
	}
	return m
}()

func (d HomeDistrict) IsValid() bool {
	_, ok := validHomeDistricts[d]
	return ok
}

func IsValidHomeDistrict(homeDistrict string) bool {
	return HomeDistrict(homeDistrict).IsValid()
}

func ParsePreferredHomeDistricts(input []string) ([]HomeDistrict, error) {
	for _, s := range input {
		if s == string(HomeDistrictAny) {
			return []HomeDistrict{HomeDistrictAny}, nil
		}
	}

	var out []HomeDistrict
	for _, s := range input {
		d := HomeDistrict(s)
		if !d.IsValid() {
			return nil, fmt.Errorf("invalid home district: %q", s)
		}
		out = append(out, d)
	}
	return out, nil
}

// used for DB storage
func HomeDistrictsToStrings(in []HomeDistrict) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = string(v)
	}
	return out
}