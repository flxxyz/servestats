package utils

import "bytes"

func TrimLine(buf *bytes.Buffer) (line []byte, err error) {
	line, err = buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(line, "\r\n"), nil
}
