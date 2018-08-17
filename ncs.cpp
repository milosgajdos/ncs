#include "ncs.h"
#include <stdio.h>

int ncs_DeviceCreate(int idx, void** deviceHandle) {
    ncStatus_t s = ncDeviceCreate(idx, (struct ncDeviceHandle_t**) deviceHandle);
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

int ncs_DeviceDestroy(void** deviceHandle) {
    ncStatus_t s = ncDeviceDestroy((struct ncDeviceHandle_t**) deviceHandle);
    return int(s);
}

int ncs_GraphCreate(const char* name, void** graphHandle) {
        ncStatus_t s = ncGraphCreate(name, (struct ncGraphHandle_t**) graphHandle);
        return int(s);
}

int ncs_GraphDestroy(void** graphHandle) {
        ncStatus_t s = ncGraphDestroy((struct ncGraphHandle_t**) graphHandle);
        return int(s);
}

int ncs_FifoCreate(const char* name, ncFifoType_t type, void** fifoHandle) {
        ncStatus_t s = ncFifoCreate(name, type, (struct ncFifoHandle_t**) fifoHandle);
        return int(s);
}

int ncs_FifoDestroy(void** fifoHandle) {
        ncStatus_t s = ncFifoDestroy((struct ncFifoHandle_t**) fifoHandle);
        return int(s);
}
