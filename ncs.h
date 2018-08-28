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
int ncs_DeviceGetOption(void* deviceHandle, int option, void *data, unsigned int *dataLength);
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
int ncs_GraphQueueInference(void* graphHandle,
                void** inFifoHandle, unsigned int inFifoCount,
                void** outFifoHandle, unsigned int outFifoCount);
int ncs_GraphQueueInferenceWithFifoElem(void* graphHandle, void* inFifoHandle, void* outFifoHandle,
                const void* inputTensor, unsigned int* inputTensorLength, void* userParam);
int ncs_GraphGetOption(void* graphHandle, int option, void *data, unsigned int *dataLength);
int ncs_GraphDestroy(void **graphHandle);

// FIFO functions
int ncs_FifoCreate(const char* name, ncFifoType_t type, void** fifoHandle);
int ncs_FifoAllocate(void* fifoHandle, void* deviceHandle, struct ncTensorDescriptor_t* tensorDesc, unsigned int numElem);

int ncs_FifoGetOption(void* fifoHandle, int option, void *data, unsigned int *dataLength);
int ncs_FifoWriteElem(void* fifoHandle, const void* inputTensor, unsigned int* inputTensorLength, void* userParam);
int ncs_FifoReadElem(void* fifoHandle, void *outputData, unsigned int* outputDataLen, void **userParam);
int ncs_FifoDestroy(void** fifoHandle);

#ifdef __cplusplus
}
#endif

#endif //_NCS_H_
