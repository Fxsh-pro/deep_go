package main

import "unsafe"

func ToLittleEndian(number uint32) uint32 {
	bytes := (*[4]byte)(unsafe.Pointer(&number))

	for i := range 2 {
		j := 3 - i
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return number
}

func toLittleEndianShifts(number uint32) uint32 {
	var res uint32
	for i := range 4 {
		res |= (number & 0xFF) << (24 - 8*i)
		number >>= 8
	}
	return res
}
