package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
)

type HomeDistrictEnum string

const (
	HomeDistrictThiruvananthapuram HomeDistrictEnum = validation.ConstantHomeDistrictThiruvananthapuram
	HomeDistrictKollam             HomeDistrictEnum = validation.ConstantHomeDistrictKollam
	HomeDistrictPathanamthitta     HomeDistrictEnum = validation.ConstantHomeDistrictPathanamthitta
	HomeDistrictAlappuzha          HomeDistrictEnum = validation.ConstantHomeDistrictAlappuzha
	HomeDistrictKottayam           HomeDistrictEnum = validation.ConstantHomeDistrictKottayam
	HomeDistrictErnakulam          HomeDistrictEnum = validation.ConstantHomeDistrictErnakulam
	HomeDistrictThrissur           HomeDistrictEnum = validation.ConstantHomeDistrictThrissur
	HomeDistrictPalakkad           HomeDistrictEnum = validation.ConstantHomeDistrictPalakkad
	HomeDistrictMalappuram         HomeDistrictEnum = validation.ConstantHomeDistrictMalappuram
	HomeDistrictKozhikode          HomeDistrictEnum = validation.ConstantHomeDistrictKozhikode
	HomeDistrictWayanad            HomeDistrictEnum = validation.ConstantHomeDistrictWayanad
	HomeDistrictKannur             HomeDistrictEnum = validation.ConstantHomeDistrictKannur
	HomeDistrictKasaragod          HomeDistrictEnum = validation.ConstantHomeDistrictKasaragod
	HomeDistrictIdukki             HomeDistrictEnum = validation.ConstantHomeDistrictIdukki
	HomeDistrictNotMentioned       HomeDistrictEnum = validation.ConstantHomeDistrictNotMentioned
)

func isValidHomeDistrictEnum(s string) bool {
	switch HomeDistrictEnum(s) {
	case HomeDistrictThiruvananthapuram, HomeDistrictKollam, HomeDistrictPathanamthitta,
		HomeDistrictAlappuzha, HomeDistrictKottayam, HomeDistrictErnakulam, HomeDistrictThrissur,
		HomeDistrictPalakkad, HomeDistrictMalappuram, HomeDistrictKozhikode, HomeDistrictWayanad,
		HomeDistrictKannur, HomeDistrictKasaragod, HomeDistrictIdukki, HomeDistrictNotMentioned:
		return true
	default:
		return false
	}
}

func (h HomeDistrictEnum) MarshalJSON() ([]byte, error) {
	if !isValidHomeDistrictEnum(string(h)) {
		return nil, fmt.Errorf("invalid home district: %q", h)
	}
	return json.Marshal(string(h))
}

func (h *HomeDistrictEnum) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !isValidHomeDistrictEnum(s) {
		return fmt.Errorf("invalid home district: %q", s)
	}
	*h = HomeDistrictEnum(s)
	return nil
}

func (h HomeDistrictEnum) Value() (driver.Value, error) {
	return string(h), nil
}

func (h *HomeDistrictEnum) Scan(value interface{}) error {
	if value == nil {
		*h = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into HomeDistrictEnum", value)
	}
	if !isValidHomeDistrictEnum(s) {
		return fmt.Errorf("invalid value for HomeDistrictEnum: %q", s)
	}
	*h = HomeDistrictEnum(s)
	return nil
}
