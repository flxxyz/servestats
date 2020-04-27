package system

import (
    "ServerStatus/cmd"
    "ServerStatus/config"
    "ServerStatus/utils"
    "fmt"
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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type SystemInfo struct {
    convert         bool
    CpuInfos        []cpu.InfoStat
    Load            *load.AvgStat
    PhysicalCpuCore int     `json:"physical_cpu_core"`
    LogicalCpuCore  int     `json:"logical_cpu_core"`
    CpuCore         int     `json:"cpu_core"`
    CpuPercent      string  `json:"cpu_percent"`
    LoadAvg         float64 `json:"load_avg"`
    MemTotal        uint64  `json:"mem_total"`
    MemTotalStr     string  `json:"mem_total_str"`
    MemFree         uint64  `json:"mem_free"`
    MemFreeStr      string  `json:"mem_free_str"`
    MemUsed         uint64  `json:"mem_used"`
    MemUsedStr      string  `json:"mem_used_str"`
    SwapTotal       uint64  `json:"swap_total"`
    SwapTotalStr    string  `json:"swap_total_str"`
    SwapFree        uint64  `json:"swap_free"`
    SwapFreeStr     string  `json:"swap_free_str"`
    SwapUsed        uint64  `json:"swap_used"`
    SwapUsedStr     string  `json:"swap_used_str"`
    HDDTotal        uint64  `json:"hdd_total"`
    HDDTotalStr     string  `json:"hdd_total_str"`
    HDDFree         uint64  `json:"hdd_free"`
    HDDFreeStr      string  `json:"hdd_free_str"`
    HDDUsed         uint64  `json:"hdd_used"`
    HDDUsedStr      string  `json:"hdd_used_str"`
    BootTime        time.Time
    Uptime          time.Time
    UptimeStr       string `json:"uptime"`
    TrafficRx       uint64 `json:"traffic_rx"`
    TrafficTx       uint64 `json:"traffic_tx"`
    TrafficRxTotal  uint64 `json:"traffic_rx_total"`
    TrafficTxTotal  uint64 `json:"traffic_tx_total"`
}

func NewSystemInfo(convert bool) *SystemInfo {
    return &SystemInfo{
        convert: convert,
    }
}

func (sys *SystemInfo) getCpuInfo() {
    infos, err := cpu.Info()
    if err != nil {
        log.Fatal("[getCpuInfo] cpu.Info() Error: ", err.Error())
    }

    for _, info := range infos {
        sys.CpuCore += int(info.Cores)
    }

    sys.CpuInfos = infos
    sys.PhysicalCpuCore, _ = cpu.Counts(false)
    sys.LogicalCpuCore, _ = cpu.Counts(true)
}

func (sys *SystemInfo) getCpuPercent() {
    var err error
    p, err := cpu.Percent(time.Second, false)
    if err != nil {
        log.Fatal("[getCpuPercent] cpu.Percent() Error: ", err.Error())
    }

    sys.CpuPercent = fmt.Sprintf("%.2f", p[0])
}

func (sys *SystemInfo) getLoad() {
    if runtime.GOOS != "windows" {
        l, err := load.Avg()
        if err != nil {
            log.Fatal("[getCpuPercent] load.Avg() Error: ", err.Error())
        }

        sys.Load = l
        sys.LoadAvg = l.Load1 / float64(sys.CpuCore)
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
    ftypes := []string{"ext4", "ext3", "ext2", "fat32", "ntfs", "xfs", "zfs", "jfs", "btrfs", "simfs", "reiserfs", "fuseblk"}

    parts, err := disk.Partitions(true)
    if err != nil {
        log.Fatal("[getHDD] disk.Partitions() Error: ", err.Error())
    }

    for _, part := range parts {
        d, err := disk.Usage(part.Mountpoint)
        if err != nil {
            log.Fatal("[getHDD] disk.Usage() Error: ", err.Error())
        }

        sort.Strings(ftypes)
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

    sys.BootTime = time.Unix(int64(timestamp), 0).Local()
}

func (sys *SystemInfo) getUptime() {
    sys.Uptime = time.Now()
    diff := sys.Uptime.Unix() - sys.BootTime.Unix()
    _, sys.UptimeStr = utils.ComputeTimeDiff(diff)
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

func (sys *SystemInfo) ConvertStr() {
    sys.MemTotalStr = utils.FileSize(sys.MemTotal)
    sys.MemFreeStr = utils.FileSize(sys.MemFree)
    sys.MemUsedStr = utils.FileSize(sys.MemUsed)

    sys.SwapTotalStr = utils.FileSize(sys.SwapTotal)
    sys.SwapFreeStr = utils.FileSize(sys.SwapFree)
    sys.SwapUsedStr = utils.FileSize(sys.SwapUsed)

    sys.HDDTotalStr = utils.FileSize(sys.HDDTotal)
    sys.HDDFreeStr = utils.FileSize(sys.HDDFree)
    sys.HDDUsedStr = utils.FileSize(sys.HDDUsed)
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

    //取网卡每秒的流量
    go sys.GetTraffic()

    if sys.convert {
        sys.ConvertStr()
    }
}

func (sys *SystemInfo) Json() (data []byte, err error) {
    data, err = json.Marshal(sys)
    return
}

func (sys *SystemInfo) JsonFormat(prefix, indent string) (data []byte, err error) {
    data, err = json.MarshalIndent(sys, prefix, indent)
    return
}

func Run(p *cmd.Cmd) {
    sys := NewSystemInfo(p.Convert)
    sys.Update()
    data, _ := sys.JsonFormat("", "  ")
    fmt.Println(string(data))
}
