package validation

const (
	ConstantCommunitySunni        = "sunni"
	ConstantCommunityMujahid      = "mujahid"
	ConstantCommunityTabligh      = "tabligh"
	ConstantCommunityJamateIslami = "jamate_islami"
	ConstantCommunityShia         = "shia"
	ConstantCommunityMuslim       = "muslim"
	ConstantCommunityNotMentioned = "not_mentioned"
)

func IsValidCommunity(community string) bool {
	switch community {
	case ConstantCommunitySunni, ConstantCommunityMujahid, ConstantCommunityTabligh,
		ConstantCommunityJamateIslami, ConstantCommunityShia,
		ConstantCommunityMuslim, ConstantCommunityNotMentioned:
		return true
	default:
		return false
	}
}
