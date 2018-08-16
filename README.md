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
	fmt.Println("Creating NCS device handle")
	dev, err := ncs.NewDevice(0)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("NCS device handle created")

	fmt.Println("Opening NCS device")
	err = dev.Open()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("NCS device opened")

	fmt.Println("Closing NCS device")
	err = dev.Close()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("NCS device closed")

	fmt.Println("Destroyig NCS device")
	err = dev.Destroy()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("NCS device destroyed")
}
```

If your Movidius NCS device is plugged in you should see the following output when running the program above:

```console
Creating NCS device handle
NCS device handle created
Opening NCS device
NCS device opened
Closing NCS device
NCS device closed
Destroyig NCS device
NCS device destroyed
```
