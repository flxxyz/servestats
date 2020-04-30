package utils

import (
	"bytes"
	"log"
	"strings"
)

func TrimLine(buf *bytes.Buffer) (line []byte, err error) {
	line, err = buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(line, "\r\n"), nil
}

func echo(m []string, v ...interface{}) []string {
	for i, _ := range v {
		switch v[i].(type) {
		case string:
			m = append(m, v[i].(string))
		case bool:
			if v[i].(bool) {
				m = append(m, "true")
			} else {
				m = append(m, "false")
			}
		case []byte:
			m = append(m, Bytes2Str(v[i].([]byte)))
		case []interface{}:
			m = echo(m, v[i].([]interface{})...)
		}
	}

	return m
}

func Echo(v ...interface{}) {
	m := append(make([]string, 0), v[0].(string))
	log.Println(strings.Join(echo(m, v[1:]...), " "))
}
