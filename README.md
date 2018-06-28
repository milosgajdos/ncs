# ncs
Neural Compute Stick V2 API Go binding

So far it only contains Device API implementation.

**WARNING, NCSDK API V2 IS BADLY BROKEN AT THE MOMENT**

# Quick start

On MacOS, clone macos branch:

```shell
$ git clone -b macos https://github.com/milosgajdos83/ncsdk.git
```

Build ncsdk API libraries:

```shell
$ cd api/src && sudo make basicinstall pythoninstall
```

Test NCSDK example:

```shell
$ cd ../../examples/apps/hello_ncs_cpp/ && make run
```

# Example Go program

```go
package main

import (
	"fmt"
	"os"

	"github.com/milosgajdos83/ncs"
)

func main() {
	ret, dev := ncs.CreateDevice(0)
	if ret != ncs.StatusOK {
		fmt.Printf("NCS Error Create(): %#v\n", ret)
		os.Exit(1)
	}

	ret = dev.Open()
	if ret != ncs.StatusOK {
		fmt.Printf("NCS Error Open(): %#v\n", ret)
		os.Exit(1)
	}

	ret = dev.Close()
	if ret != ncs.StatusOK {
		fmt.Printf("NCS Error Close(): %#v\n", ret)
		os.Exit(1)
	}

	ret = ncs.DestroyDevice(dev)
	if ret != ncs.StatusOK {
		fmt.Printf("NCS Error Destroy(): %#v\n", ret)
		os.Exit(1)
	}
}
```

This will most likely print the same errors as the official C++ example released by Intel - the API V2 is badly broken as of now:

```
Can't create semaphore
: Function not implemented
Can't create semaphore
: Function not implemented
E: [         0] dispatcherAddEvent:533	can't wait semaphore

W: [         0] dispatcherAddEvent:545	No more semaphores. Increase XLink or OS resources

E: [         0] dispatcherAddEvent:533	can't wait semaphore

W: [         0] dispatcherAddEvent:545	No more semaphores. Increase XLink or OS resources

E: [         0] XLinkOpenStream:909	Max streamId reached deaddead!
W: [         0] ncDeviceOpen:558	can't open stream

NCS Error Open(): -2
```

