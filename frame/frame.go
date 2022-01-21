package frame

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type _ = strings.Builder
type _ = unsafe.Pointer

var _ = math.Float32frombits
var _ = math.Float64frombits
var _ = strconv.FormatInt
var _ = strconv.FormatUint
var _ = strconv.FormatFloat
var _ = fmt.Sprint

type Frame []byte

func (s Frame) Service() uint32 {
	_ = s[3]
	var __v uint32 = uint32(s[0]) |
		uint32(s[1])<<8 |
		uint32(s[2])<<16 |
		uint32(s[3])<<24
	return uint32(__v)
}

func (s Frame) Method() uint32 {
	_ = s[7]
	var __v uint32 = uint32(s[4]) |
		uint32(s[5])<<8 |
		uint32(s[6])<<16 |
		uint32(s[7])<<24
	return uint32(__v)
}

func (s Frame) Sequence() uint32 {
	_ = s[11]
	var __v uint32 = uint32(s[8]) |
		uint32(s[9])<<8 |
		uint32(s[10])<<16 |
		uint32(s[11])<<24
	return uint32(__v)
}

func (s Frame) BodySize() uint32 {
	_ = s[15]
	var __v uint32 = uint32(s[12]) |
		uint32(s[13])<<8 |
		uint32(s[14])<<16 |
		uint32(s[15])<<24
	return uint32(__v)
}

func (s Frame) Body() []byte {
	_ = s[23]
	var __off0 uint64 = 24
	var __off1 uint64 = uint64(s[16]) |
		uint64(s[17])<<8 |
		uint64(s[18])<<16 |
		uint64(s[19])<<24 |
		uint64(s[20])<<32 |
		uint64(s[21])<<40 |
		uint64(s[22])<<48 |
		uint64(s[23])<<56
	return []byte(s[__off0:__off1])
}

func (s Frame) Vstruct_Validate() bool {
	if len(s) < 24 {
		return false
	}

	var __off0 uint64 = 24
	var __off1 uint64 = uint64(s[16]) |
		uint64(s[17])<<8 |
		uint64(s[18])<<16 |
		uint64(s[19])<<24 |
		uint64(s[20])<<32 |
		uint64(s[21])<<40 |
		uint64(s[22])<<48 |
		uint64(s[23])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func (s Frame) String() string {
	if !s.Vstruct_Validate() {
		return "Frame (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Frame {")
	__b.WriteString("Service: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Service()), 10))
	__b.WriteString(", ")
	__b.WriteString("Method: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Method()), 10))
	__b.WriteString(", ")
	__b.WriteString("Sequence: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Sequence()), 10))
	__b.WriteString(", ")
	__b.WriteString("BodySize: ")
	__b.WriteString(strconv.FormatUint(uint64(s.BodySize()), 10))
	__b.WriteString(", ")
	__b.WriteString("Body: ")
	__b.WriteString(fmt.Sprint(s.Body()))
	__b.WriteString("}")
	return __b.String()
}

func Serialize_Frame(dst Frame, Service uint32, Method uint32, Sequence uint32, BodySize uint32, Body []byte) Frame {
	_ = dst[23]
	dst[0] = byte(Service)
	dst[1] = byte(Service >> 8)
	dst[2] = byte(Service >> 16)
	dst[3] = byte(Service >> 24)
	dst[4] = byte(Method)
	dst[5] = byte(Method >> 8)
	dst[6] = byte(Method >> 16)
	dst[7] = byte(Method >> 24)
	dst[8] = byte(Sequence)
	dst[9] = byte(Sequence >> 8)
	dst[10] = byte(Sequence >> 16)
	dst[11] = byte(Sequence >> 24)
	dst[12] = byte(BodySize)
	dst[13] = byte(BodySize >> 8)
	dst[14] = byte(BodySize >> 16)
	dst[15] = byte(BodySize >> 24)

	var __index = uint64(24)
	__tmp_4 := uint64(len(Body)) + __index
	dst[16] = byte(__tmp_4)
	dst[17] = byte(__tmp_4 >> 8)
	dst[18] = byte(__tmp_4 >> 16)
	dst[19] = byte(__tmp_4 >> 24)
	dst[20] = byte(__tmp_4 >> 32)
	dst[21] = byte(__tmp_4 >> 40)
	dst[22] = byte(__tmp_4 >> 48)
	dst[23] = byte(__tmp_4 >> 56)
	copy(dst[__index:__tmp_4], Body)
	return dst
}

func New_Frame(Service uint32, Method uint32, Sequence uint32, BodySize uint32, Body []byte) Frame {
	var __vstruct__size = 24 + len(Body)
	var __vstruct__buf = make(Frame, __vstruct__size)
	__vstruct__buf = Serialize_Frame(__vstruct__buf, Service, Method, Sequence, BodySize, Body)
	return __vstruct__buf
}
