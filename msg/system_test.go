package msg

import (
	"bytes"
	"github.com/flxxyz/ServerStatus/config"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"testing"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TestOS(t *testing.T) {
	os := NewOS()

	switch runtime.GOOS {
	case "windows":
		versions := map[string]string{
			`5\.0`:  "2000",
			`5\.1`:  "XP",
			`5\.2`:  "Server 2003",
			`6\.0`:  "Vista",
			`6\.1`:  "7",
			`6\.2`:  "8",
			`6\.3`:  "8.1",
			`10\.0`: "10",
		}

		cmd := exec.Command("cmd.exe")
		out, _ := cmd.StdoutPipe()
		buffer := bytes.NewBuffer(make([]byte, 0))
		cmd.Start()
		buffer.ReadFrom(out)
		str, _ := buffer.ReadString(']')
		cmd.Wait()
		for key, _ := range versions {
			re := regexp.MustCompile(`Microsoft Windows \[[\s\S]* ` + key + `\.[0-9]+\.[0-9]+\]`)
			if re.MatchString(str) {
				os.Version = versions[key]
				//return os
				break
			}
		}
	case "linux":
		if ok, _ := PathExists("/etc/os-release"); ok {
			cmd := exec.Command("cat", "/etc/os-release")
			stdout, _ := cmd.StdoutPipe()
			cmd.Start()
			content, err := ioutil.ReadAll(stdout)
			if err == nil {
				id := regexp.MustCompile(`ID="?(.*?)"?\n`).FindStringSubmatch(string(content))
				if len(id) > 1 {
					os.Name = id[1]
				}

				versionId := regexp.MustCompile(`VERSION_ID="?([.0-9]+)"?\n`).FindStringSubmatch(string(content))
				if len(versionId) > 1 {
					os.Version = versionId[1]
				}
			}
		}
	}

	t.Log(os)
}

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
