#ifndef _NCS_H_
#define _NCS_H_

#include <stdlib.h>
#include <mvnc.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef ncFifoType_t ncFifoType;
typedef ncFifoDataType_t ncFifoDataType;

int ncs_DeviceCreate(int idx, void **deviceHandle);
int ncs_DeviceOpen(void* deviceHandle);
int ncs_DeviceClose(void* deviceHandle);
int ncs_DeviceDestroy(void **deviceHandle);

int ncs_GraphCreate(const char* name, void **graphHandle);
int ncs_GraphDestroy(void **graphHandle);

int ncs_FifoCreate(const char* name, ncFifoType_t type, void** fifoHandle);
int ncs_FifoDestroy(void** fifoHandle);

#ifdef __cplusplus
}
#endif

#endif //_NCS_H_
