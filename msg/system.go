package msg

import (
	"log"
	"os/exec"
	"runtime"
	"sort"
	"time"

	"github.com/flxxyz/ServerStatus/config"
	"github.com/flxxyz/ServerStatus/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var (
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
)

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

func NewUnknownOS() *OS {
	return &OS{
		Name:    "unknown",
		Version: "unknown",
		Arch:    "unknown",
	}
}

type SystemInfo struct {
	err             error
	cpuInfos        []cpu.InfoStat
	load            *load.AvgStat
	PhysicalCpuCore uint64     `json:"physical_cpu_core" xml:"physical_cpu_core"`
	LogicalCpuCore  uint64     `json:"logical_cpu_core" xml:"logical_cpu_core"`
	CpuCore         uint64     `json:"cpu_core" xml:"cpu_core"`
	CpuPercent      float64    `json:"cpu_percent" xml:"cpu_percent"`
	LoadAvg         [3]float64 `json:"load_avg" xml:"load_avg"`
	MemTotal        uint64     `json:"mem_total" xml:"mem_total"`
	MemFree         uint64     `json:"mem_free" xml:"mem_free"`
	MemUsed         uint64     `json:"mem_used" xml:"mem_used"`
	SwapTotal       uint64     `json:"swap_total" xml:"swap_total"`
	SwapFree        uint64     `json:"swap_free" xml:"swap_free"`
	SwapUsed        uint64     `json:"swap_used" xml:"swap_used"`
	HDDTotal        uint64     `json:"hdd_total" xml:"hdd_total"`
	HDDFree         uint64     `json:"hdd_free" xml:"hdd_free"`
	HDDUsed         uint64     `json:"hdd_used" xml:"hdd_used"`
	BootTime        time.Time  `json:"-" xml:"-"`
	Uptime          uint64     `json:"uptime" xml:"uptime"`
	TrafficRx       uint64     `json:"traffic_rx" xml:"traffic_rx"`
	TrafficRxTotal  uint64     `json:"traffic_rx_total" xml:"traffic_rx_total"`
	TrafficTx       uint64     `json:"traffic_tx" xml:"traffic_tx"`
	TrafficTxTotal  uint64     `json:"traffic_tx_total" xml:"traffic_tx_total"`
	IPv4Support     bool       `json:"ipv4_support" xml:"ipv4_support"`
	IPv6Support     bool       `json:"ipv6_support" xml:"ipv6_support"`
	OS              *OS        `json:"os" xml:"os"`
}

func NewSystemInfo() (sys *SystemInfo) {
	sys = &SystemInfo{
		OS: NewOS(),
	}
	switch runtime.GOOS {
	case "windows":
		sys.OS.Version = utils.GetWindowsVersion()
	case "linux":
		sys.OS.Name, sys.OS.Version = utils.GetLinuxVersion()
	}
	return
}

func (sys *SystemInfo) getCpuInfo() {
	sys.cpuInfos, sys.err = cpu.Info()
	if sys.err != nil {
		log.Println("[getCpuInfo] cpu.Info() Error: ", sys.err.Error())
		return
	}

	physicalCore, _ := cpu.Counts(false)
	logicalCore, _ := cpu.Counts(true)
	sys.PhysicalCpuCore = uint64(physicalCore)
	sys.LogicalCpuCore = uint64(logicalCore)

	sys.CpuCore = 0
	for _, info := range sys.cpuInfos {
		sys.CpuCore += uint64(info.Cores)
	}

	if sys.LogicalCpuCore > sys.PhysicalCpuCore {
		sys.CpuCore = sys.LogicalCpuCore
	} else {
		if sys.PhysicalCpuCore > sys.LogicalCpuCore {
			sys.CpuCore = sys.PhysicalCpuCore
		}
	}
}

func (sys *SystemInfo) getCpuPercent() {
	p, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("[getCpuPercent] cpu.Percent() Error: ", err.Error())
		return
	}

	sys.CpuPercent = utils.Decimal(p[0], 2)
}

func (sys *SystemInfo) getLoad() {
	if runtime.GOOS != "windows" {
		sys.load, sys.err = load.Avg()
		if sys.err != nil {
			log.Println("[getLoad] load.Avg() Error: ", sys.err.Error())
			return
		}

		sys.LoadAvg[0] = utils.Decimal(sys.load.Load1/float64(sys.CpuCore), 2)
		sys.LoadAvg[1] = utils.Decimal(sys.load.Load5/float64(sys.CpuCore), 2)
		sys.LoadAvg[2] = utils.Decimal(sys.load.Load15/float64(sys.CpuCore), 2)
	}
}

func (sys *SystemInfo) getMem() {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Println("[getMem] mem.VirtualMemory() Error: ", err.Error())
		return
	}

	sys.MemTotal = m.Total
	sys.MemFree = m.Free
	sys.MemUsed = m.Used
}

func (sys *SystemInfo) getSwap() {
	s, err := mem.SwapMemory()
	if err != nil {
		log.Println("[getSwap] mem.SwapMemory() Error: ", err.Error())
		return
	}

	sys.SwapTotal = s.Total
	sys.SwapFree = s.Free
	sys.SwapUsed = s.Used
}

func (sys *SystemInfo) getHDD() {
	parts, err := disk.Partitions(true)
	if err != nil {
		log.Println("[getHDD] disk.Partitions() Error: ", err.Error())
		return
	}

	sys.resetHDD()

	for _, part := range parts {
		d, err := disk.Usage(part.Mountpoint)
		if err != nil {
			log.Println("[getHDD] disk.Usage() Error: ", err.Error())
			return
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
		log.Println("[getBootTime] host.BootTime() Error: ", err.Error())
		return
	}

	sys.BootTime = time.Unix(int64(timestamp), 0).Local()
}

func (sys *SystemInfo) getUptime() {
	sys.Uptime = uint64(time.Now().Unix() - sys.BootTime.Unix())
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

		time.Sleep(config.IvGetTraffic)
	}
}

func (sys *SystemInfo) CheckIPvNSupport() {
	ipv4Command := "ping"
	ipv4CommandArgs := []string{"-4", config.DnsIpv4, "-c", "3"}
	ipv6Command := "ping"
	ipv6CommandArgs := []string{"-6", config.DnsIpv6, "-c", "3"}
	if runtime.GOOS == "windows" {
		ipv4CommandArgs[2] = "-n"
		ipv6CommandArgs[2] = "-n"
	} else {
		if runtime.GOOS == "darwin" {
			ipv4CommandArgs = append(ipv4CommandArgs[:0], ipv4CommandArgs[1:]...)
			ipv6Command = "ping6"
			ipv6CommandArgs = append(ipv6CommandArgs[:0], ipv6CommandArgs[1:]...)
		}
	}

	for {
		cmdIPv4 := exec.Command(ipv4Command, ipv4CommandArgs...)
		if cmdIPv4.Run() != nil {
			sys.IPv4Support = false
		} else {
			sys.IPv4Support = true
		}

		cmdIPv6 := exec.Command(ipv6Command, ipv6CommandArgs...)
		if cmdIPv6.Run() != nil {
			sys.IPv6Support = false
		} else {
			sys.IPv6Support = true
		}

		time.Sleep(config.IvCheckIPvNSupport)
	}
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
}

func (sys *SystemInfo) Reset() {
	sys.PhysicalCpuCore = 0
	sys.LogicalCpuCore = 0
	sys.CpuCore = 0
	sys.CpuPercent = 0
	sys.LoadAvg = [3]float64{0, 0, 0}
	sys.MemTotal = 0
	sys.MemFree = 0
	sys.MemUsed = 0
	sys.SwapTotal = 0
	sys.SwapFree = 0
	sys.SwapUsed = 0
	sys.resetHDD()
	sys.Uptime = 0
	sys.TrafficRx = 0
	sys.TrafficRxTotal = 0
	sys.TrafficTx = 0
	sys.TrafficTxTotal = 0
	sys.IPv4Support = false
	sys.IPv6Support = false
}

func (sys *SystemInfo) resetHDD() {
	sys.HDDTotal = 0
	sys.HDDFree = 0
	sys.HDDUsed = 0
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

func init() {
	sort.Strings(ftypes)
}
