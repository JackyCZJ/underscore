package openwrt

/*
#include <stdlib.h>
#include <string.h>
#include <libubox/blob.h>
#include <libubox/blobmsg_json.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"

	"git.esixcloud.net/flash/underscore/json"
)

type BlobBuf struct {
	ptr *C.struct_blob_buf
}

func NewBlobBuf() *BlobBuf {
	var p = C.calloc(1, C.sizeof_struct_blob_buf)
	var _ = C.memset(p, 0, C.sizeof_struct_blob_buf)
	var B = &BlobBuf{
		ptr: (*C.struct_blob_buf)(p),
	}
	runtime.SetFinalizer(B, (*BlobBuf).Free)
	return B
}

func (buf *BlobBuf) Init(id int) int {
	if buf == nil || buf.ptr == nil {
		return -1
	}
	return int(C.blob_buf_init(buf.ptr, C.int(id)))
}

func (buf *BlobBuf) Free() {
	if buf == nil || buf.ptr == nil {
		return
	}
	runtime.SetFinalizer(buf, nil)
	C.blob_buf_free(buf.ptr)
	C.free(unsafe.Pointer(buf.ptr))
	buf.ptr = nil
}

func (buf *BlobBuf) AddJsonFrom(obj any) error {
	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case string:
		if err := buf.AddJsonFromString(v); err != nil {
			return err
		}
	default:
		ret, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		if err := buf.AddJsonFromString(string(ret)); err != nil {
			return err
		}
	}

	return nil
}

func (buf *BlobBuf) AddJsonFromString(str string) error {
	if buf == nil || buf.ptr == nil {
		return errors.New("blob buffer has been freed")
	}
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	_, err := C.blobmsg_add_json_from_string(buf.ptr, cstr)
	return err
}
