package core

// VisionType is a creature's relevant sense for a situation. VisionNormal is the
// baseline every creature has.
type VisionType int

const (
	VisionNormal VisionType = iota
	VisionDarkvision
	VisionBlindsight
	VisionTruesight
	VisionTremorsense
)

// LightLevel is the ambient light. LightUnspecified is the must-set zero.
type LightLevel int

const (
	LightUnspecified LightLevel = iota
	LightBright
	LightDim
	LightDark
)

// Visibility is the result of SightVisibility. VisibilityClear is the safe zero.
type Visibility int

const (
	VisibilityClear Visibility = iota
	VisibilityObscured
	VisibilityBlocked
)

// EffectiveLight applies a creature's vision to the ambient light. Darkvision in
// range treats dark as dim and dim as bright; blindsight and truesight see
// regardless of light.
func EffectiveLight(ambient LightLevel, vision VisionType, withinRange bool) LightLevel {
	switch vision {
	case VisionBlindsight, VisionTruesight, VisionTremorsense:
		if withinRange {
			return LightBright
		}
		return ambient
	case VisionDarkvision:
		if !withinRange {
			return ambient
		}
		switch ambient {
		case LightDark:
			return LightDim
		case LightDim:
			return LightBright
		default:
			return ambient
		}
	default:
		return ambient
	}
}

// SightVisibility turns effective light and conditions into a visibility state:
// dim is Obscured (disadvantage on sight Perception); darkness, blinded, or an
// invisible target is Blocked (cannot see).
func SightVisibility(light LightLevel, blinded, targetInvisible bool) Visibility {
	if blinded || targetInvisible || light == LightDark {
		return VisibilityBlocked
	}
	if light == LightDim {
		return VisibilityObscured
	}
	return VisibilityClear
}
