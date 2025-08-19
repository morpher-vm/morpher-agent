package collector

import (
	"net"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
)

type OSInfo struct {
	HostName      string `json:"hostname"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	KernelVersion string `json:"kernel_version"`
}

type CPUInfo struct {
	Arch      string  `json:"arch"`
	VCPUs     int     `json:"vcpus"`
	Model     string  `json:"model"`
	MHzPerCPU float64 `json:"mhz_per_cpu"`
}

type RAMInfo struct {
	TotalMB uint64 `json:"total_mb"`
	UsedMB  uint64 `json:"used_mb"`
	FreeMB  uint64 `json:"free_mb"`
}

type DiskInfo struct {
	Mounts []DiskMount `json:"mounts"`
}

type DiskMount struct {
	Mount       string  `json:"mount"`
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type NetworkInterface struct {
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardware_addr"`
	MTU          int      `json:"mtu"`
	Addrs        []string `json:"addrs"`
}

type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

type SystemInfo struct {
	ID      *string      `json:"id"`
	OS      *OSInfo      `json:"os"`
	CPU     *CPUInfo     `json:"cpu"`
	RAM     *RAMInfo     `json:"ram"`
	Disk    *DiskInfo    `json:"disk"`
	Network *NetworkInfo `json:"network"`
}

func CollectOS() (*OSInfo, error) {
	var out OSInfo

	if hi, err := host.Info(); err == nil {
		out.HostName = hi.Hostname
		out.Name = titleOrRaw(hi.Platform)
		out.Version = hi.PlatformVersion
		out.KernelVersion = hi.KernelVersion
	}

	return &out, nil
}

func CollectCPU() (*CPUInfo, error) {
	var out CPUInfo

	out.Arch = runtime.GOARCH
	if vcpus, err := cpu.Counts(true); err == nil {
		out.VCPUs = vcpus
	}
	if infos, err := cpu.Info(); err == nil && len(infos) > 0 {
		out.Model = infos[0].ModelName
		var sum float64
		var n float64
		for _, it := range infos {
			if it.Mhz > 0 {
				sum += it.Mhz
				n++
			}
		}
		if n > 0 {
			out.MHzPerCPU = round1(sum / n)
		}
	}

	return &out, nil
}

func CollectRAM() (*RAMInfo, error) {
	var out RAMInfo

	if vm, err := mem.VirtualMemory(); err == nil {
		out.TotalMB = vm.Total / 1024 / 1024
		used := (vm.Total - vm.Available) / 1024 / 1024
		out.UsedMB = used
		out.FreeMB = vm.Available / 1024 / 1024
	}

	return &out, nil
}

func CollectDisk() (*DiskInfo, error) {
	var out DiskInfo

	out.Mounts = make([]DiskMount, 0, 8)
	if parts, err := disk.Partitions(true); err == nil {
		for _, p := range parts {
			if skipFS(p.Fstype) || skipMount(p.Mountpoint) {
				continue
			}
			if u, err := disk.Usage(p.Mountpoint); err == nil && u.Total > 0 {
				out.Mounts = append(out.Mounts, DiskMount{
					Mount:       p.Mountpoint,
					TotalGB:     round1(bytesToGB(u.Total)),
					UsedGB:      round1(bytesToGB(u.Used)),
					UsedPercent: round1(u.UsedPercent),
				})
			}
		}
	}

	return &out, nil
}

func CollectNetwork() (*NetworkInfo, error) {
	var out NetworkInfo

	if ifaces, err := gnet.Interfaces(); err == nil {
		out.Interfaces = make([]NetworkInterface, 0, len(ifaces))
		for _, iface := range ifaces {
			// Skip interfaces that are down or loopback.
			has := func(w string) bool {
				for _, f := range iface.Flags {
					if strings.EqualFold(f, w) {
						return true
					}
				}
				return false
			}
			if !has("up") || has("loopback") {
				continue
			}

			var addrs []string
			for _, a := range iface.Addrs {
				ip := a.Addr
				if i := strings.IndexByte(ip, '%'); i != -1 {
					ip = ip[:i]
				}
				if j := strings.IndexByte(ip, '/'); j != -1 {
					ip = ip[:j]
				}

				p := net.ParseIP(ip)
				if ip == "" || p == nil || p.IsLoopback() || p.IsMulticast() ||
					p.IsUnspecified() || p.IsLinkLocalUnicast() || p.IsLinkLocalMulticast() {
					continue
				}
				addrs = append(addrs, ip)
			}
			if len(addrs) == 0 {
				continue
			}
			out.Interfaces = append(out.Interfaces, NetworkInterface{
				Name:         iface.Name,
				HardwareAddr: iface.HardwareAddr,
				MTU:          iface.MTU,
				Addrs:        addrs,
			})
		}
	}
	return &out, nil
}

func CollectBiosId() (*string, error) {
	id, err := host.HostID()
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func CollectSystem() (*SystemInfo, error) {
	os, err := CollectOS()
	if err != nil {
		return nil, err
	}

	cpu, err := CollectCPU()
	if err != nil {
		return nil, err
	}

	ram, err := CollectRAM()
	if err != nil {
		return nil, err
	}

	disk, err := CollectDisk()
	if err != nil {
		return nil, err
	}

	network, err := CollectNetwork()
	if err != nil {
		return nil, err
	}

	biosId, err := CollectBiosId()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		ID:      biosId,
		OS:      os,
		CPU:     cpu,
		RAM:     ram,
		Disk:    disk,
		Network: network,
	}, nil
}

func bytesToGB(b uint64) float64 { return float64(b) / (1024 * 1024 * 1024) }

func round1(x float64) float64 {
	if x < 0 {
		return float64(int64(x*10-0.5)) / 10
	}
	return float64(int64(x*10+0.5)) / 10
}

func titleOrRaw(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func skipFS(fs string) bool {
	fs = strings.ToLower(fs)
	bad := []string{
		"tmpfs", "devtmpfs", "proc", "sysfs", "overlay", "squashfs",
		"autofs", "cgroup", "cgroup2", "pstore", "tracefs", "debugfs",
		"devfs", "aufs", "ramfs", "fusectl", "mqueue", "bpf",
	}
	for _, v := range bad {
		if fs == v {
			return true
		}
	}
	return false
}

func skipMount(m string) bool {
	prefixes := []string{"/proc", "/sys", "/run", "/dev", "/var/lib/docker", "/var/lib/containers"}
	for _, p := range prefixes {
		if strings.HasPrefix(m, p) {
			return true
		}
	}
	return false
}
