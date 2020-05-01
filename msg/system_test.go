package msg

import (
	"github.com/flxxyz/ServerStatus/config"
	"os/exec"
	"runtime"
	"testing"
)

func TestIPvN(t *testing.T) {
	if runtime.GOOS == "windows" {
		count = "-n"
		timeout = "-w"
	}

	cmdIPv4 := exec.Command("ping", "-4", config.DnsIpv4, count, "1", timeout, "5")
	errIPv4 := cmdIPv4.Run()
	if errIPv4 == nil {
		t.Log("IPv4 Support")
	} else {
		t.Log("IPv4 Not Support")
	}

	cmdIPv6 := exec.Command("ping", "-6", config.DnsIpv6, "-c", "1", "-W", "5")
	errIPv6 := cmdIPv6.Run()
	if errIPv6 == nil {
		t.Log("IPv6 Support")
	} else {
		t.Log("IPv6 Not Support")
	}
}
