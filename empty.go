package vrpc

type EmptyValue []byte

func Empty() EmptyValue {
	return newEmpty([]byte{})
}

func (s EmptyValue) Vstruct_Validate() bool {
	if len(s) < 8 {
		return false
	}

	var __off0 uint64 = 8
	var __off1 uint64 = uint64(s[0]) |
		uint64(s[1])<<8 |
		uint64(s[2])<<16 |
		uint64(s[3])<<24 |
		uint64(s[4])<<32 |
		uint64(s[5])<<40 |
		uint64(s[6])<<48 |
		uint64(s[7])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func serializeEmpty(dst EmptyValue, Any []byte) EmptyValue {
	_ = dst[7]

	var __index = uint64(8)
	__tmp_0 := uint64(len(Any)) + __index
	dst[0] = byte(__tmp_0)
	dst[1] = byte(__tmp_0 >> 8)
	dst[2] = byte(__tmp_0 >> 16)
	dst[3] = byte(__tmp_0 >> 24)
	dst[4] = byte(__tmp_0 >> 32)
	dst[5] = byte(__tmp_0 >> 40)
	dst[6] = byte(__tmp_0 >> 48)
	dst[7] = byte(__tmp_0 >> 56)
	copy(dst[__index:__tmp_0], Any)
	return dst
}

func newEmpty(Any []byte) EmptyValue {
	var __vstruct__size = 8 + len(Any)
	var __vstruct__buf = make(EmptyValue, __vstruct__size)
	__vstruct__buf = serializeEmpty(__vstruct__buf, Any)
	return __vstruct__buf
}
