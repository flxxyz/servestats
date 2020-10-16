package msg

import (
	"testing"
)

func TestIPvN(t *testing.T) {
	s := NewSystemInfo()
	s.CheckIPvNSupport()

	if s.IPv4Support {
		t.Log("IPv4 Support")
	} else {
		t.Log("IPv4 Not Support")
	}

	if s.IPv6Support {
		t.Log("IPv6 Support")
	} else {
		t.Log("IPv6 Not Support")
	}
}
