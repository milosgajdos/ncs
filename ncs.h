#ifndef _NCS_H_
#define _NCS_H_

#include <stdlib.h>
#include <mvnc.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef ncFifoType_t ncFifoType;
typedef ncFifoDataType_t ncFifoDataType;
// Device Functions
int ncs_DeviceCreate(int idx, void **deviceHandle);
int ncs_DeviceOpen(void* deviceHandle);
int ncs_DeviceClose(void* deviceHandle);
int ncs_DeviceDestroy(void **deviceHandle);
// Graph Functions
int ncs_GraphCreate(const char* name, void **graphHandle);
int ncs_GraphAllocate(void* deviceHandle, void* graphHandle,
                const void *graphBuffer, unsigned int graphBufferLength);
int ncs_GraphAllocateWithFifos(void* deviceHandle, void* graphHandle,
                const void *graphBuffer, unsigned int graphBufferLength,
                void** inFifoHandle, void** outFifoHandle);
int ncs_GraphAllocateWithFifosEx(void* deviceHandle, void* graphHandle,
                const void *graphBuffer, unsigned int graphBufferLength,
                void** inFifoHandle, ncFifoType_t inFifoType, int inNumElem, ncFifoDataType_t inDataType,
                void** outFifoHandle, ncFifoType_t outFifoType, int outNumElem, ncFifoDataType_t outDataType);
int ncs_GraphDestroy(void **graphHandle);
// FIFO functions
int ncs_FifoCreate(const char* name, ncFifoType_t type, void** fifoHandle);
int ncs_FifoDestroy(void** fifoHandle);

#ifdef __cplusplus
}
#endif

#endif //_NCS_H_
