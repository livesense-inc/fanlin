package lcms

// #cgo CFLAGS: -std=c11 -D_POSIX_C_SOURCE=200809 -Wall -Wextra -Wpedantic -Wundef -O3
// #cgo LDFLAGS: -llcms2
// #include <lcms2.h>
import "C"
