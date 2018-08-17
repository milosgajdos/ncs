package ncs

// #cgo LDFLAGS: -lmvnc
/*
#include <ncs.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// StatusCode is the NCSDK API status code as returned by most API calls.
// It usually reports the status of the neural compute stick.
type StatusCode int

const (
	// StatusOK means the API function call worked as expected
	StatusOK StatusCode = 0

	// StatusBusy means device is busy, retry later.
	StatusBusy = -1

	// StatusError means an unexpected error was encountered during the API function call.
	StatusError = -2

	// StatusOutOfMemory means the host is out of memory.
	StatusOutOfMemory = -3

	// StatusDeviceNotFound means no device has been found at the given index or name.
	StatusDeviceNotFound = -4

	// StatusInvalidParameters means at least one of the given parameters is wrong.
	StatusInvalidParameters = -5

	// StatusTimeout means there was a timeout in the communication with the device.
	StatusTimeout = -6

	// StatusCmdNotFound means the file to boot the device was not found.
	StatusCmdNotFound = -7

	// StatusNotAllocated means the graph or fifo has not been allocated.
	StatusNotAllocated = -8

	// StatusUnauthorized means an unauthorized operation has been attempted.
	StatusUnauthorized = -9

	// StatusUnsupportedGraphFile means the graph file version is not supported.
	StatusUnsupportedGraphFile = -10

	// StatusUnsupportedConfigFile is reserved for future use.
	StatusUnsupportedConfigFile = -11

	// StatusUnsupportedFeature means the operation used a feature unsupported by this firmware version.
	StatusUnsupportedFeature = -12

	// StatusMyriadError when an error has been reported by device, use MVNC_DEBUG_INFO.
	StatusMyriadError = -13

	// StatusInvalidDataLength means an invalid data length has been passed when getting or setting an option
	StatusInvalidDataLength = -14

	// StatusInvalidHandle means an invalid handle has been passed to a function
	StatusInvalidHandle = -15
)

// String method to satisfy fmt.Stringer interface
func (nc StatusCode) String() string {
	switch nc {
	case StatusOK:
		return "OK"
	case StatusBusy:
		return "Device busy"
	case StatusError:
		return "Unexpected error"
	case StatusOutOfMemory:
		return "Host out of memory"
	case StatusDeviceNotFound:
		return "Device not found"
	case StatusInvalidParameters:
		return "Invalid parameters"
	case StatusTimeout:
		return "Device timeout"
	case StatusCmdNotFound:
		return "Device bootloader not found"
	case StatusNotAllocated:
		return "Unallocated resource"
	case StatusUnauthorized:
		return "Unauthorized operation"
	case StatusUnsupportedGraphFile:
		return "Unsupported graph file"
	case StatusUnsupportedConfigFile:
		return "Unsupported configuration"
	case StatusUnsupportedFeature:
		return "Unsupported feature"
	case StatusMyriadError:
		return "Movidius VPU failure"
	case StatusInvalidDataLength:
		return "Invalid data length when setting options"
	case StatusInvalidHandle:
		return "Invalid handle"
	default:
		return "Unknown status"
	}
}

// Device is Neural Compute Stick (NCS) device
type Device struct {
	handle unsafe.Pointer
}

// NewDevice creates new NCS device handle and returns it
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceCreate.html
func NewDevice(index int) (*Device, error) {
	var handle unsafe.Pointer

	c := C.ncs_DeviceCreate(C.int(index), &handle)

	if StatusCode(c) == StatusOK {
		return &Device{handle: handle}, nil
	}

	return nil, fmt.Errorf("Failed to create new device: %s", StatusCode(c))
}

// Open initializes NCS device and opens device communication channel
// It returns error if it fails to open the device.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceOpen.html
func (d *Device) Open() error {
	c := C.ncs_DeviceOpen(d.handle)

	if StatusCode(c) == StatusOK {
		return nil
	}

	return fmt.Errorf("Failed to open device: %s", StatusCode(c))
}

// Close closes the communication channel with NCS device.
// It returns error if it fails to close the communication channel.
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceClose.html
func (d *Device) Close() error {
	c := C.ncs_DeviceClose(d.handle)

	if StatusCode(c) == StatusOK {
		return nil
	}

	return fmt.Errorf("Failed to close device: %s", StatusCode(c))
}

// Destroy destroys NCS device handle and frees associated resources.
// This function must be called for every device that was initialized with NewDevice().
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncDeviceDestroy.html
func (d *Device) Destroy() error {
	c := C.ncs_DeviceDestroy(&d.handle)

	if StatusCode(c) == StatusOK {
		return nil
	}

	return fmt.Errorf("Failed to destroy device: %s", StatusCode(c))
}

// Graph is NCSDK neural network graph
type Graph struct {
	handle unsafe.Pointer
	d      *Device
}

// NewGraph creates new Graph with given name and returns it
// It returns error if it fails to create new graph
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphCreate.html
func NewGraph(name string) (*Graph, error) {
	var handle unsafe.Pointer

	_name := C.CString(name)
	defer C.free(unsafe.Pointer(_name))

	c := C.ncs_GraphCreate(_name, &handle)

	if StatusCode(c) == StatusOK {
		return &Graph{handle: handle}, nil
	}

	return nil, fmt.Errorf("Failed to create new graph: %s", StatusCode(c))
}

// Destroy destroys NCS graph handle and frees associated resources.
// This function must be called for every graph that was initialized with NewGraph().
//
// For more information:
// https://movidius.github.io/ncsdk/ncapi/ncapi2/c_api/ncGraphDestroy.html
func (g *Graph) Destroy() error {
	c := C.ncs_GraphDestroy(&g.handle)

	if StatusCode(c) == StatusOK {
		return nil
	}

	return fmt.Errorf("Failed to destroy graph: %s", StatusCode(c))
}
