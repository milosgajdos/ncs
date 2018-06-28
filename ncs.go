package ncs

// #cgo LDFLAGS: -lmvnc
/*
#include <ncs.h>
*/
import "C"
import "unsafe"

// Status is the NCSDK API status code returned by most API calls.
// It usually reports the status of the neural compute stick.
type Status int

const (
	// StatusOK when the API function call worked as expected
	StatusOK Status = 0

	// StatusBusy means device is busy, retry later.
	StatusBusy = -1

	// StatusError means an unexpected error was encountered during the API function call.
	StatusError = -2

	// StatusOutOfMemory means the host is out of memory.
	StatusOutOfMemory = -3

	// StatusDeviceNotFound means no device at the given index or name.
	StatusDeviceNotFound = -4

	// StatusInvalidParameters when at least one of the given parameters is wrong.
	StatusInvalidParameters = -5

	// StatusTimeout means there was a timeout in the communication with the device.
	StatusTimeout = -6

	// StatusCmdNotFound means the file to boot the device was not found.
	StatusCmdNotFound = -7

	// StatusNotAllocated means the graph or fifo has not been allocated
	StatusNotAllocated = -8

	// StatusUnauthorized means an unauthorized operation has been attempted
	StatusUnauthorized = -9

	// StatusUnsupportedGraphFile means the graph file version is not supported.
	StatusUnsupportedGraphFile = -10

	// StatusUnsupportedConfigFile is reserved for future use
	StatusUnsupportedConfigFile = -11

	// StatusUnsupportedFeature means the operation used a feature unsupported by this firmware version
	StatusUnsupportedFeature = -12

	// StatusMyriadError when an error has been reported by the device, use MVNC_DEBUG_INFO.
	StatusMyriadError = -13

	// StatusInvalidDataLength means an invalid data length has been passed when getting or setting an option
	StatusInvalidDataLength = -14

	// StatusInvalidHandle means an n invalid handle has been passed to a function
	StatusInvalidHandle = -15
)

// LogLevel defines application logging levels
type LogLevel int

const (
	// Debug logs debug and above (full verbosity)
	Debug LogLevel = iota

	// Info logs info and above
	Info

	// Warn logs warning and above (default)
	Warn

	// Error logs errors and above
	Error

	// Fatal logs fatal errors only
	Fatal
)

// Device is Neural Compute Stick (NCS) device
type Device struct {
	handle unsafe.Pointer
}

// CreateDevice creates new NCS device and returns it
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceCreate.html
func CreateDevice(index int) (Status, *Device) {
	var handle unsafe.Pointer
	s := C.ncs_DeviceCreate(C.int(index), &handle)

	return Status(s), &Device{handle: handle}
}

// Open initializes a NCS device and opens communication
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceOpen.html
func (d *Device) Open() Status {
	ret := C.ncs_DeviceOpen(d.handle)

	return Status(ret)
}

// Close closes communication with a NCS device
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceClose.html
func (d *Device) Close() Status {
	res := C.ncs_DeviceClose(d.handle)

	return Status(res)
}

// DestroyDevice destroys a handle for a NCS compute device and frees associated resources.
// This function must be called for every device that was initialized with CreateDevice()
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceDestroy.html
func DestroyDevice(d *Device) Status {
	s := C.ncs_DeviceDestroy(&d.handle)

	return Status(s)
}
