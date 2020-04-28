package msg

import (
	"ServerStatus/config"
	"ServerStatus/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
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
)

type SystemInfo struct {
	err               error
	cpuInfos          []cpu.InfoStat
	load              *load.AvgStat
	PhysicalCpuCore   int        `json:"physical_cpu_core"`
	LogicalCpuCore    int        `json:"logical_cpu_core"`
	CpuCore           int        `json:"cpu_core"`
	CpuPercent        float64    `json:"cpu_percent"`
	LoadAvg           [3]float64 `json:"load_avg"`
	MemTotal          uint64     `json:"mem_total"`
	MemFree           uint64     `json:"mem_free"`
	MemUsed           uint64     `json:"mem_used"`
	MemTotalStr       string     `json:"mem_total_str"`
	MemFreeStr        string     `json:"mem_free_str"`
	MemUsedStr        string     `json:"mem_used_str"`
	SwapTotal         uint64     `json:"swap_total"`
	SwapFree          uint64     `json:"swap_free"`
	SwapUsed          uint64     `json:"swap_used"`
	SwapTotalStr      string     `json:"swap_total_str"`
	SwapFreeStr       string     `json:"swap_free_str"`
	SwapUsedStr       string     `json:"swap_used_str"`
	HDDTotal          uint64     `json:"hdd_total"`
	HDDFree           uint64     `json:"hdd_free"`
	HDDUsed           uint64     `json:"hdd_used"`
	HDDTotalStr       string     `json:"hdd_total_str"`
	HDDFreeStr        string     `json:"hdd_free_str"`
	HDDUsedStr        string     `json:"hdd_used_str"`
	bootTime          time.Time
	uptime            time.Time
	Uptime            int64  `json:"uptime"`
	UptimeStr         string `json:"uptime_str"`
	TrafficRx         uint64 `json:"traffic_rx"`
	TrafficRxStr      string `json:"traffic_rx_str"`
	TrafficRxTotal    uint64 `json:"traffic_rx_total"`
	TrafficRxTotalStr string `json:"traffic_rx_total_str"`
	TrafficTx         uint64 `json:"traffic_tx"`
	TrafficTxStr      string `json:"traffic_tx_str"`
	TrafficTxTotal    uint64 `json:"traffic_tx_total"`
	TrafficTxTotalStr string `json:"traffic_tx_total_str"`
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
	diff := sys.uptime.Unix() - sys.bootTime.Unix()
	sys.Uptime, sys.UptimeStr = utils.ComputeTimeDiff(diff)
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

		time.Sleep(time.Second * config.IntervalRefreshTraffic)
	}
}

func (sys *SystemInfo) toString() {
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

	sys.toString()
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

func init() {
	sort.Strings(ftypes)
}
