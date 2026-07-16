package onvif

import (
	"github.com/nguyencobap/onvif/xsd/onvif"
)

// Relative pan/tilt translation spaces, from annex A of the PTZ specification. A
// camera lists the ones it takes per node, and what a value means depends entirely
// on which is named: the same 0.5 is half a picture in one and half a full sweep of
// the head in the other.
const (
	// PanTiltSpaceFov normalizes against the current field of view: the corner of
	// the picture is (1,1) whatever the lens is doing, so a point on the picture is
	// a translation without anything else being known about the camera. Optional.
	PanTiltSpaceFov = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/TranslationSpaceFov"
	// PanTiltSpaceGeneric normalizes against the whole range of the head instead, so
	// the same value covers far more ground zoomed in than zoomed out. Mandatory.
	PanTiltSpaceGeneric = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/TranslationGenericSpace"
)

// RelativePanTiltSpace returns the relative pan/tilt space with this URI as the
// camera described it when it was reached, or nil if it never named it.
//
// Callers should treat not naming a space as not implementing it, because asking is
// the only way to find out: a camera sent a space it does not implement does not
// refuse it, it quietly moves in its own default instead.
func (dev *Device) RelativePanTiltSpace(uri string) *onvif.Space2DDescription {
	for _, node := range dev.ptzNodes {
		if node.SupportedPTZSpaces == nil {
			continue
		}
		spaces := node.SupportedPTZSpaces.RelativePanTiltTranslationSpace
		for i, space := range spaces {
			if space.URI != nil && string(*space.URI) == uri {
				return &spaces[i]
			}
		}
	}
	return nil
}
