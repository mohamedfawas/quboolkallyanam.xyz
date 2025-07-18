package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type HomeDistrictEnum string

const (
	HomeDistrictThiruvananthapuram HomeDistrictEnum = "thiruvananthapuram"
	HomeDistrictKollam             HomeDistrictEnum = "kollam"
	HomeDistrictPathanamthitta     HomeDistrictEnum = "pathanamthitta"
	HomeDistrictAlappuzha          HomeDistrictEnum = "alappuzha"
	HomeDistrictKottayam           HomeDistrictEnum = "kottayam"
	HomeDistrictErnakulam          HomeDistrictEnum = "ernakulam"
	HomeDistrictThrissur           HomeDistrictEnum = "thrissur"
	HomeDistrictPalakkad           HomeDistrictEnum = "palakkad"
	HomeDistrictMalappuram         HomeDistrictEnum = "malappuram"
	HomeDistrictKozhikode          HomeDistrictEnum = "kozhikode"
	HomeDistrictWayanad            HomeDistrictEnum = "wayanad"
	HomeDistrictKannur             HomeDistrictEnum = "kannur"
	HomeDistrictKasaragod          HomeDistrictEnum = "kasaragod"
	HomeDistrictIdukki             HomeDistrictEnum = "idukki"
	HomeDistrictNotMentioned       HomeDistrictEnum = "not_mentioned"
)

func (h HomeDistrictEnum) Value() (driver.Value, error) {
	return string(h), nil
}

func (h *HomeDistrictEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok {
		*h = HomeDistrictEnum(s)
		return nil
	}
	return fmt.Errorf("cannot scan %T into HomeDistrictEnum", value)
}

func (p *PartnerPreference) SetHomeDistricts(districts []HomeDistrictEnum) error {
	stringDistricts := make([]string, len(districts))
	for i, d := range districts {
		stringDistricts[i] = string(d)
	}
	data, err := json.Marshal(stringDistricts)
	if err != nil {
		return err
	}
	p.PreferredHomeDistricts = datatypes.JSON(data)
	return nil
}

func (p *PartnerPreference) GetHomeDistricts() ([]HomeDistrictEnum, error) {
	var stringDistricts []string
	if len(p.PreferredHomeDistricts) == 0 {
		return []HomeDistrictEnum{}, nil
	}

	if err := json.Unmarshal(p.PreferredHomeDistricts, &stringDistricts); err != nil {
		return nil, err
	}

	districts := make([]HomeDistrictEnum, len(stringDistricts))
	for i, s := range stringDistricts {
		districts[i] = HomeDistrictEnum(s)
	}
	return districts, nil
}
