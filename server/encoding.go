package server

import (
	"bytes"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	littleEndian = "LE"
	bigEndian    = "BE"
)

func utf16ToString(data []byte) string {
	utf16ByteOrder := isUTF16(data)

	if utf16ByteOrder == "" {
		return string(data)
	}

	u16s := make([]uint16, 1)
	buffer := &bytes.Buffer{}
	utf8Buffer := make([]byte, 4)

	for i := 0; i < len(data); i += 2 {
		if utf16ByteOrder == littleEndian {
			u16s[0] = uint16(data[i]) + (uint16(data[i+1]) << 8)
		} else {
			u16s[0] = uint16(data[i+1]) + (uint16(data[i]) << 8)
		}

		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(utf8Buffer, r[0])
		buffer.Write(utf8Buffer[:n])
	}

	return buffer.String()
}

func isUTF16(data []byte) string {
	if len(data) >= 2 {
		if data[0] == 0xFE && data[1] == 0xFF {
			return bigEndian
		} else if data[0] == 0xFF && data[1] == 0xFE {
			return littleEndian
		}
	}

	return ""
}
