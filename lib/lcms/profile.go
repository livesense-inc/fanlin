package lcms

// #include <stdlib.h>
// #include <lcms2.h>
import "C"

import (
	"unsafe"
)

type Profile struct {
	prof C.cmsHPROFILE
}

func (prof *Profile) CloseProfile() {
	if prof.prof != nil {
		C.cmsCloseProfile(prof.prof)
	}
}

func OpenProfileFromMem(profdata []byte) *Profile {
	data := unsafe.Pointer(&profdata[0])
	dataLen := C.cmsUInt32Number(len(profdata))
	return &Profile{prof: C.cmsOpenProfileFromMem(data, dataLen)}
}

func CreateSRGBProfile() *Profile {
	return &Profile{prof: C.cmsCreate_sRGBProfile()}
}
