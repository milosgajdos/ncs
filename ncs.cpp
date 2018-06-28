#include "ncs.h"
#include <stdio.h>

int ncs_DeviceCreate(int idx, void** deviceHandle) {
    ncStatus_t s = ncDeviceCreate(idx, (struct ncDeviceHandle_t**) deviceHandle);
    return int(s);
}

int ncs_DeviceDestroy(void** deviceHandle) {
    ncStatus_t s = ncDeviceDestroy((struct ncDeviceHandle_t**) deviceHandle);
    return int(s);
}

int ncs_DeviceOpen(void* deviceHandle) {
    ncStatus_t s = ncDeviceOpen((struct ncDeviceHandle_t*) deviceHandle);
    return int(s);
}

int ncs_DeviceClose(void* deviceHandle) {
    ncStatus_t s = ncDeviceClose((struct ncDeviceHandle_t*) deviceHandle);
    return int(s);
}
