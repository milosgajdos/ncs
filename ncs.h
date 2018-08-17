#ifndef _NCS_H_
#define _NCS_H_

#include <stdlib.h>
#include <mvnc.h>

#ifdef __cplusplus
extern "C" {
#endif

int ncs_DeviceCreate(int idx, void **deviceHandle);
int ncs_DeviceOpen(void* deviceHandle);
int ncs_DeviceClose(void* deviceHandle);
int ncs_DeviceDestroy(void **deviceHandle);

int ncs_GraphCreate(const char* name, void **graphHandle);
int ncs_GraphDestroy(void **graphHandle);

#ifdef __cplusplus
}
#endif

#endif //_NCS_H_
