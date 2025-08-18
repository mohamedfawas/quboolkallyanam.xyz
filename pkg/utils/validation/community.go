package validation

import "fmt"

type Community string

const (
	Sunni        Community = "sunni"
	Mujahid      Community = "mujahid"
	Tabligh      Community = "tabligh"
	JamateIslami Community = "jamate_islami"
	Shia         Community = "shia"
	Muslim       Community = "muslim"
	CommunityAny Community = "any"
)

var validSet = func() map[Community]struct{} {
	vals := []Community{
		Sunni,
		Mujahid,
		Tabligh,
		JamateIslami,
		Shia,
		Muslim,
		CommunityAny,
	}
	m := make(map[Community]struct{}, len(vals))
	for _, c := range vals {
		m[c] = struct{}{}
	}
	return m
}()

func (c Community) IsValid() bool {
	_, ok := validSet[c]
	return ok
}

func IsValidCommunity(community string) bool {
	return Community(community).IsValid()
}

func ParsePreferredCommunities(input []string) ([]Community, error) {
    // If "any" was requested, ignore everything else and just return [CommunityAny].
    for _, s := range input {
        if s == string(CommunityAny) {
            return []Community{CommunityAny}, nil
        }
    }

    // Otherwise validate and collect every community.
    var out []Community
    for _, s := range input {
        c := Community(s)
        if !c.IsValid() {
            return nil, fmt.Errorf("invalid community: %q", s)
        }
        out = append(out, c)
    }
    return out, nil
}

// used for DB storage
func CommunitiesToStrings(in []Community) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = string(v)
	}
	return out
}
