package utils

import (
	"bytes"
	"github.com/flxxyz/ServerStatus/msg"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

func TrimLine(buf *bytes.Buffer) (line []byte, err error) {
	line, err = buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(line, "\r\n"), nil
}

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

func GetWindowsVersion() (version string) {
	version = "unknown"
	versions := map[string]string{
		`5\.0`:  "Windows 2000",
		`5\.1`:  "Windows XP",
		`5\.2`:  "Windows Server 2003",
		`6\.0`:  "Windows Vista",
		`6\.1`:  "Windows 7",
		`6\.2`:  "Windows 8",
		`6\.3`:  "Windows 8.1",
		`10\.0`: "Windows 10",
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
			return versions[key]
		}
	}

	return
}

func GetLinuxVersion() *msg.OS {
	name := runtime.GOOS
	version := "unknown"
	if ok, _ := PathExists("/etc/os-release"); ok {
		cmd := exec.Command("cat", "/etc/os-release")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		content, err := ioutil.ReadAll(stdout)
		if err == nil {
			id := regexp.MustCompile(`ID=(.*?)\n`).FindStringSubmatch(string(content))
			if len(id) > 1 {
				name = id[1]
			}

			versionId := regexp.MustCompile(`VERSION_ID="?([.0-9]+)"?\n`).FindStringSubmatch(string(content))
			if len(versionId) > 1 {
				version = versionId[1]
			}
		}
	}

	return &msg.OS{
		Name:    name,
		Version: version,
		Arch:    runtime.GOARCH,
	}
}
