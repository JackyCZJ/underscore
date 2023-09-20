package openwrt

/*
#include <libubox/blob.h>
#include <libubox/blobmsg_json.h>
*/
import "C"
import (
	"sync"
	"unsafe"

	"git.esixcloud.net/flash/underscore/json"
)

type BlobBuf struct {
	ptr *C.struct_blob_buf
}

var p = C.calloc(1, C.sizeof_struct_blob_buf)
var _ = C.memset(p, 0, C.sizeof_struct_blob_buf)
var B = &BlobBuf{
	ptr: (*C.struct_blob_buf)(p),
}
var lock sync.Mutex

func NewBlobBuf() *BlobBuf {
	return B
}

func (buf *BlobBuf) Init(id int) int {
	lock.Lock()
	return int(C.blob_buf_init(buf.ptr, C.int(id)))
}

func (buf *BlobBuf) Free() {
	lock.Unlock()
	C.blob_buf_free(buf.ptr)
}

func (buf *BlobBuf) AddJsonFrom(obj any) error {
	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case string:
		buf.AddJsonFromString(v)
	default:
		ret, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		buf.AddJsonFromString(string(ret))
	}

	return nil
}

func (buf *BlobBuf) AddJsonFromString(str string) error {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	_, err := C.blobmsg_add_json_from_string(buf.ptr, cstr)
	return err
}
