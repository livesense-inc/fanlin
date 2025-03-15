package lcms

// #include <lcms2.h>
import "C"

type CMSType int

const (
	TYPE_RGBA_8 CMSType = C.TYPE_RGBA_8
	TYPE_CMYK_8 CMSType = C.TYPE_CMYK_8
)
