package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type CommunityEnum string

const (
	CommunitySunni        CommunityEnum = "sunni"
	CommunityMujahid      CommunityEnum = "mujahid"
	CommunityTabligh      CommunityEnum = "tabligh"
	CommunityJamateIslami CommunityEnum = "jamate_islami"
	CommunityShia         CommunityEnum = "shia"
	CommunityMuslim       CommunityEnum = "muslim"
	CommunityNotMentioned CommunityEnum = "not_mentioned"
)

func (c CommunityEnum) Value() (driver.Value, error) {
	return string(c), nil
}

func (c *CommunityEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok {
		*c = CommunityEnum(s)
		return nil
	}
	return fmt.Errorf("cannot scan %T into CommunityEnum", value)
}

// Helper methods for JSONB operations on PartnerPreference
func (p *PartnerPreference) SetCommunities(communities []CommunityEnum) error {
	stringCommunities := make([]string, len(communities))
	for i, c := range communities {
		stringCommunities[i] = string(c)
	}
	data, err := json.Marshal(stringCommunities)
	if err != nil {
		return err
	}
	p.PreferredCommunities = datatypes.JSON(data)
	return nil
}

func (p *PartnerPreference) GetCommunities() ([]CommunityEnum, error) {
	var stringCommunities []string
	if len(p.PreferredCommunities) == 0 {
		return []CommunityEnum{}, nil
	}

	if err := json.Unmarshal(p.PreferredCommunities, &stringCommunities); err != nil {
		return nil, err
	}

	communities := make([]CommunityEnum, len(stringCommunities))
	for i, s := range stringCommunities {
		communities[i] = CommunityEnum(s)
	}
	return communities, nil
}
