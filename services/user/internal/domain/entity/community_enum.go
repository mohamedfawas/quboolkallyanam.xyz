package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

func isValidCommunityEnum(val string) bool {
	switch CommunityEnum(val) {
	case CommunitySunni, CommunityMujahid, CommunityTabligh,
		CommunityJamateIslami, CommunityShia,
		CommunityMuslim, CommunityNotMentioned:
		return true
	default:
		return false
	}
}

// - Ensures the enum is valid before converting it to a JSON string.
//
// Example use case:
//
//	json.Marshal(CommunitySunni) → `"sunni"`
//	json.Marshal("invalid") → error
func (c CommunityEnum) MarshalJSON() ([]byte, error) {
	if !isValidCommunityEnum(string(c)) {
		return nil, fmt.Errorf("invalid community: %q", c)
	}
	return json.Marshal(string(c))
}

// - Parses the string from JSON and checks if it's a valid enum.
// Example use case:
//
//	json.Unmarshal([]byte(`"shia"`), &c) → c = CommunityShia
//	json.Unmarshal([]byte(`"invalid"`), &c) → error
func (c *CommunityEnum) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !isValidCommunityEnum(s) {
		return fmt.Errorf("invalid community: %q", s)
	}
	*c = CommunityEnum(s)
	return nil
}

// Value is used when saving the enum to the database (e.g., using GORM).
// Example use case:
//
//	INSERT INTO partner_preferences (preferred_communities) VALUES ('["sunni", "muslim"]')
func (c CommunityEnum) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan is used when reading the enum value from the database into Go.
// Example use case:
//
//	DB has JSONB: ["sunni", "mujahid"]
//	→ GORM reads it, calls Scan("sunni") → sets Go field
func (c *CommunityEnum) Scan(src interface{}) error {
	if src == nil {
		*c = ""
		return nil
	}
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into CommunityEnum", src)
	}
	if !isValidCommunityEnum(s) {
		return fmt.Errorf("invalid community: %s", s)
	}
	*c = CommunityEnum(s)
	return nil
}
