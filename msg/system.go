package msg

import (
	"bytes"
	"github.com/flxxyz/ServerStatus/config"
	"github.com/flxxyz/ServerStatus/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"time"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	ftypes = []string{
		"ext4", //貌似现在不知此ext4，会识别成ext2/ext3
		"ext3",
		"ext2",
		"ext2/ext3",
		"ext",
		"fat",
		"msdos",
		"ntfs",
		"xfs",
		"zfs",
		"jfs",
		"btrfs",
		"simfs",
		"reiserfs",
		"fuseblk",
	}
	count   = "-c"
	timeout = "-W"
)

type SystemInfo struct {
	hasConvStr        bool
	err               error
	cpuInfos          []cpu.InfoStat
	load              *load.AvgStat
	PhysicalCpuCore   int        `json:"physical_cpu_core" xml:"physical_cpu_core"`
	LogicalCpuCore    int        `json:"logical_cpu_core" xml:"logical_cpu_core"`
	CpuCore           int        `json:"cpu_core" xml:"cpu_core"`
	CpuPercent        float64    `json:"cpu_percent" xml:"cpu_percent"`
	LoadAvg           [3]float64 `json:"load_avg" xml:"load_avg"`
	MemTotal          uint64     `json:"mem_total" xml:"mem_total"`
	MemFree           uint64     `json:"mem_free" xml:"mem_free"`
	MemUsed           uint64     `json:"mem_used" xml:"mem_used"`
	MemTotalStr       string     `json:"mem_total_str" xml:"mem_total_str"`
	MemFreeStr        string     `json:"mem_free_str" xml:"mem_free_str"`
	MemUsedStr        string     `json:"mem_used_str" xml:"mem_used_str"`
	SwapTotal         uint64     `json:"swap_total" xml:"swap_total"`
	SwapFree          uint64     `json:"swap_free" xml:"swap_free"`
	SwapUsed          uint64     `json:"swap_used" xml:"swap_used"`
	SwapTotalStr      string     `json:"swap_total_str" xml:"swap_total_str"`
	SwapFreeStr       string     `json:"swap_free_str" xml:"swap_free_str"`
	SwapUsedStr       string     `json:"swap_used_str" xml:"swap_used_str"`
	HDDTotal          uint64     `json:"hdd_total" xml:"hdd_total"`
	HDDFree           uint64     `json:"hdd_free" xml:"hdd_free"`
	HDDUsed           uint64     `json:"hdd_used" xml:"hdd_used"`
	HDDTotalStr       string     `json:"hdd_total_str" xml:"hdd_total_str"`
	HDDFreeStr        string     `json:"hdd_free_str" xml:"hdd_free_str"`
	HDDUsedStr        string     `json:"hdd_used_str" xml:"hdd_used_str"`
	bootTime          time.Time
	uptime            time.Time
	Uptime            int64  `json:"uptime" xml:"uptime"`
	UptimeStr         string `json:"uptime_str" xml:"uptime_str"`
	TrafficRx         uint64 `json:"traffic_rx" xml:"traffic_rx"`
	TrafficRxStr      string `json:"traffic_rx_str" xml:"traffic_rx_str"`
	TrafficRxTotal    uint64 `json:"traffic_rx_total" xml:"traffic_rx_total"`
	TrafficRxTotalStr string `json:"traffic_rx_total_str" xml:"traffic_rx_total_str"`
	TrafficTx         uint64 `json:"traffic_tx" xml:"traffic_tx"`
	TrafficTxStr      string `json:"traffic_tx_str" xml:"traffic_tx_str"`
	TrafficTxTotal    uint64 `json:"traffic_tx_total" xml:"traffic_tx_total"`
	TrafficTxTotalStr string `json:"traffic_tx_total_str" xml:"traffic_tx_total_str"`
	IPv4Support       bool   `json:"ipv4_support" xml:"ipv4_support"`
	IPv6Support       bool   `json:"ipv6_support" xml:"ipv6_support"`
	OS                *OS    `json:"os" xml:"os"`
}

func NewSystemInfo(hasConvStr bool) *SystemInfo {
	os := NewOS()
	switch runtime.GOOS {
	case "windows":
		os = GetWindowsVersion()
	case "linux":
		os = GetLinuxVersion()
	}

	return &SystemInfo{
		hasConvStr: hasConvStr,
		OS:         os,
	}
}

func (sys *SystemInfo) getCpuInfo() {
	sys.cpuInfos, sys.err = cpu.Info()
	if sys.err != nil {
		log.Fatal("[getCpuInfo] cpu.Info() Error: ", sys.err.Error())
	}

	sys.CpuCore = 0
	for _, info := range sys.cpuInfos {
		sys.CpuCore += int(info.Cores)
	}

	sys.PhysicalCpuCore, _ = cpu.Counts(false)
	sys.LogicalCpuCore, _ = cpu.Counts(true)
}

func (sys *SystemInfo) getCpuPercent() {
	p, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatal("[getCpuPercent] cpu.Percent() Error: ", err.Error())
	}

	sys.CpuPercent = utils.Decimal(p[0], 2)
}

func (sys *SystemInfo) getLoad() {
	if runtime.GOOS != "windows" {
		sys.load, sys.err = load.Avg()
		if sys.err != nil {
			log.Fatal("[getCpuPercent] load.Avg() Error: ", sys.err.Error())
		}

		sys.LoadAvg[0] = utils.Decimal(sys.load.Load1/float64(sys.CpuCore), 2)
		sys.LoadAvg[1] = utils.Decimal(sys.load.Load5/float64(sys.CpuCore), 2)
		sys.LoadAvg[2] = utils.Decimal(sys.load.Load15/float64(sys.CpuCore), 2)
	}
}

func (sys *SystemInfo) getMem() {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal("[getCpuPercent] mem.VirtualMemory() Error: ", err.Error())
	}

	sys.MemTotal = m.Total
	sys.MemFree = m.Free
	sys.MemUsed = m.Used
}

func (sys *SystemInfo) getSwap() {
	s, err := mem.SwapMemory()
	if err != nil {
		log.Fatal("[getCpuPercent] mem.SwapMemory() Error: ", err.Error())
	}

	sys.SwapTotal = s.Total
	sys.SwapFree = s.Free
	sys.SwapUsed = s.Used
}

func (sys *SystemInfo) getHDD() {
	parts, err := disk.Partitions(true)
	if err != nil {
		log.Fatal("[getHDD] disk.Partitions() Error: ", err.Error())
	}

	sys.resetHDD()

	for _, part := range parts {
		d, err := disk.Usage(part.Mountpoint)
		if err != nil {
			log.Fatal("[getHDD] disk.Usage() Error: ", err.Error())
		}

		i := sort.SearchStrings(ftypes, d.Fstype)
		if i < len(ftypes) && ftypes[i] == d.Fstype {
			sys.HDDTotal += d.Total
			sys.HDDFree += d.Free
			sys.HDDUsed += d.Used
		}
	}
}

func (sys *SystemInfo) getBootTime() {
	timestamp, err := host.BootTime()
	if err != nil {
		log.Fatal("[getBootTime] host.BootTime() Error: ", err.Error())
	}

	sys.bootTime = time.Unix(int64(timestamp), 0).Local()
}

func (sys *SystemInfo) getUptime() {
	sys.uptime = time.Now()
	sys.Uptime = sys.uptime.Unix() - sys.bootTime.Unix()
	if _, sys.UptimeStr = utils.ComputeTimeDiff(sys.Uptime); !sys.hasConvStr {
		sys.UptimeStr = ""
	}
}

func (sys *SystemInfo) GetTraffic() {
	for {
		infos, _ := net.IOCounters(false)
		netCard := infos[0]

		//下行流量
		if sys.TrafficRxTotal != 0 {
			sys.TrafficRx = netCard.BytesRecv - sys.TrafficRxTotal
		}
		sys.TrafficRxTotal = netCard.BytesRecv

		//上行流量
		if sys.TrafficTxTotal != 0 {
			sys.TrafficTx = netCard.BytesSent - sys.TrafficTxTotal
		}
		sys.TrafficTxTotal = netCard.BytesSent

		time.Sleep(config.IntervalGetTraffic)
	}
}

func (sys *SystemInfo) CheckIPvNSupport() {
	for {
		if runtime.GOOS == "windows" {
			count = "-n"
			timeout = "-w"
		}

		cmd := exec.Command("ping", "-4", config.DnsIpv4, count, "1", timeout, "5")
		err := cmd.Run()
		if err != nil {
			sys.IPv4Support = false
		} else {
			sys.IPv4Support = true
		}

		cmdIPv6 := exec.Command("ping", "-6", config.DnsIpv6, count, "1", timeout, "5")
		errIPv6 := cmdIPv6.Run()
		if errIPv6 != nil {
			sys.IPv6Support = false
		} else {
			sys.IPv6Support = true
		}

		time.Sleep(config.IntervalCheckIPvNSupport)
	}
}

func (sys *SystemInfo) ToString() {
	sys.MemTotalStr = utils.FileSize(sys.MemTotal)
	sys.MemFreeStr = utils.FileSize(sys.MemFree)
	sys.MemUsedStr = utils.FileSize(sys.MemUsed)

	sys.SwapTotalStr = utils.FileSize(sys.SwapTotal)
	sys.SwapFreeStr = utils.FileSize(sys.SwapFree)
	sys.SwapUsedStr = utils.FileSize(sys.SwapUsed)

	sys.HDDTotalStr = utils.FileSize(sys.HDDTotal)
	sys.HDDFreeStr = utils.FileSize(sys.HDDFree)
	sys.HDDUsedStr = utils.FileSize(sys.HDDUsed)

	sys.TrafficRxStr = utils.FileSize(sys.TrafficRx)
	sys.TrafficTxStr = utils.FileSize(sys.TrafficTx)
	sys.TrafficRxTotalStr = utils.FileSize(sys.TrafficRxTotal)
	sys.TrafficTxTotalStr = utils.FileSize(sys.TrafficTxTotal)
}

func (sys *SystemInfo) Update() {
	sys.getCpuInfo()
	sys.getCpuPercent()
	sys.getLoad()
	sys.getMem()
	sys.getSwap()
	sys.getHDD()
	sys.getBootTime()
	sys.getUptime()

	if sys.hasConvStr {
		sys.ToString()
	}
}

func (sys *SystemInfo) Reset() {
	sys.PhysicalCpuCore = 0
	sys.LogicalCpuCore = 0
	sys.CpuCore = 0
	sys.CpuPercent = 0
	sys.LoadAvg = [3]float64{0, 0, 0}
	sys.MemTotal = 0
	sys.MemTotalStr = ""
	sys.MemFree = 0
	sys.MemFreeStr = ""
	sys.MemUsed = 0
	sys.MemUsedStr = ""
	sys.SwapTotal = 0
	sys.SwapTotalStr = ""
	sys.SwapFree = 0
	sys.SwapFreeStr = ""
	sys.SwapUsed = 0
	sys.SwapUsedStr = ""
	sys.resetHDD()
	sys.Uptime = 0
	sys.UptimeStr = ""
	sys.TrafficRx = 0
	sys.TrafficRxStr = ""
	sys.TrafficRxTotal = 0
	sys.TrafficRxTotalStr = ""
	sys.TrafficTx = 0
	sys.TrafficTxStr = ""
	sys.TrafficTxTotal = 0
	sys.TrafficTxTotalStr = ""
	sys.IPv4Support = false
	sys.IPv6Support = false
	sys.OS = NewOS()
}

func (sys *SystemInfo) resetHDD() {
	sys.HDDTotal = 0
	sys.HDDTotalStr = ""
	sys.HDDFree = 0
	sys.HDDFreeStr = ""
	sys.HDDUsed = 0
	sys.HDDUsedStr = ""
}

func (sys *SystemInfo) Set(data []byte) {
	_ = json.Unmarshal(data, sys)
}

func (sys *SystemInfo) Json() (data []byte, err error) {
	data, err = json.Marshal(sys)
	return
}

func (sys *SystemInfo) JsonFormat(prefix, indent string) (data []byte, err error) {
	data, err = json.MarshalIndent(sys, prefix, indent)
	return
}

type OS struct {
	Name    string `json:"name" xml:"name"`
	Version string `json:"version" xml:"version"`
	Arch    string `json:"arch" xml:"arch"`
}

func NewOS() *OS {
	return &OS{
		Name:    runtime.GOOS,
		Version: "unknown",
		Arch:    runtime.GOARCH,
	}
}

func GetWindowsVersion() (os *OS) {
	os = NewOS()

	versionNumbers := map[string]string{
		`5\.0`:  "2000",
		`5\.1`:  "XP",
		`5\.2`:  "Server 2003",
		`6\.0`:  "Server 2008",
		`6\.1`:  "Server 2008 R2",
		`6\.2`:  "Server 2012",
		`6\.3`:  "Server 2012 R2",
		`10\.0`: "10",
	}

	//win10VersionNumbers := map[string]string{
	//	`10\.0\.14300`: "Server 2016",
	//	`10\.0\.14393`: "Server 2016",
	//	`10\.0\.16299`: "Server 2016",
	//	`10\.0\.17134`: "Server 2016",
	//	`10\.0\.17677`: "Server 2019",
	//	`10\.0\.17763`: "Server 2019",
	//	`10\.0\.18362`: "Server 2019",
	//	`10\.0\.18363`: "Server 2019",
	//}

	cmd := exec.Command("cmd.exe")
	out, _ := cmd.StdoutPipe()
	buffer := bytes.NewBuffer(make([]byte, 0))
	cmd.Start()
	buffer.ReadFrom(out)
	str, _ := buffer.ReadString(']')
	cmd.Wait()
	for key, _ := range versionNumbers {
		re := regexp.MustCompile(`Microsoft Windows \[[\s\S]* ` + key + `\.([0-9]+).?[0-9]*\]`)
		if re.MatchString(str) {
			if versionNumbers[key] != "10" {
				os.Version = versionNumbers[key]
			} else {
				versionNumber := re.FindStringSubmatch(str)
				if len(versionNumber) > 1 {
					if utils.Str2Int(versionNumber[1]) <= 17134 {
						os.Version = "Server 2016"
					} else {
						os.Version = "Server 2019"
					}
				}
			}

			return
		}
	}

	return
}

func GetLinuxVersion() (os *OS) {
	os = NewOS()

	if ok, _ := utils.PathExists("/etc/os-release"); ok {
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

	return
}

func init() {
	sort.Strings(ftypes)
}
