package main

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"strings"
	"path/filepath"
	"io/ioutil"
	"github.com/shirou/gopsutil/process"
)

// service info path
var service_dir = "/tmp/gameservice"

var services = make(map[int]int) 

var zoneId = 0

var done = make(chan int)

// 
// service_[id] -> PID
// 
func loadServiceInfo() {

	// load zone info
	content, err := ioutil.ReadFile(service_dir + "/zone_id")
	if err == nil {
		text := strings.TrimSpace(string(content))
		zoneId, _ = strconv.Atoi(text)
	}

	// load services's pid
	files, _:= filepath.Glob(service_dir + "/service_*")

	for i := 0; i < len(files); i++ {

		_, filename := filepath.Split(files[i]) 
		
		index := strings.Index(filename, "_")
		if index > -1 {
			service_id, err := strconv.Atoi(filename[index+1:])
			if err == nil {
				content, err := ioutil.ReadFile(files[i])
				if err == nil {
					pidtext := strings.TrimSpace(string(content))
					pid, _ := strconv.Atoi(pidtext)
					services[service_id] = pid
				}
			}
		} 
	}
}

func printServiceMetrics(service int, pid int) {
	// not exist
	{
		ret, _ := process.PidExists(int32(pid))
		if ret == false {
			fmt.Printf("service_crash,zone=%d,service=%d value=1\n", zoneId, service)
			done <- 0
			return
		}
	}

	ret, err := process.NewProcess(int32(pid))
	if err == nil {
		printMetrics(service, *ret)
	}

	done <- 1
}

func printMetrics(service int, proc process.Process) {
    var prefix string

	fields := map[string]interface{}{}

	numThreads, err := proc.NumThreads()
	if err == nil {
		fields[prefix+"num_threads"] = numThreads
    }
    
	// fds, err := proc.NumFDs()
	// if err == nil {
	// 	fields[prefix+"num_fds"] = fds
	// }

	// ctx, err := proc.NumCtxSwitches()
	// if err == nil {
	// 	fields[prefix+"voluntary_context_switches"] = ctx.Voluntary
	// 	fields[prefix+"involuntary_context_switches"] = ctx.Involuntary
	// }

	// io, err := proc.IOCounters()
	// if err == nil {
	// 	fields[prefix+"read_count"] = io.ReadCount
	// 	fields[prefix+"write_count"] = io.WriteCount
	// 	fields[prefix+"read_bytes"] = io.ReadBytes
	// 	fields[prefix+"write_bytes"] = io.WriteBytes
	// }

	// cpu_time, err := proc.Times()
	// if err == nil {
	// 	fields[prefix+"cpu_time_user"] = cpu_time.User
	// 	fields[prefix+"cpu_time_system"] = cpu_time.System
	// 	fields[prefix+"cpu_time_idle"] = cpu_time.Idle
	// 	fields[prefix+"cpu_time_nice"] = cpu_time.Nice
	// 	fields[prefix+"cpu_time_iowait"] = cpu_time.Iowait
	// 	fields[prefix+"cpu_time_irq"] = cpu_time.Irq
	// 	fields[prefix+"cpu_time_soft_irq"] = cpu_time.Softirq
	// 	fields[prefix+"cpu_time_steal"] = cpu_time.Steal
	// 	fields[prefix+"cpu_time_stolen"] = cpu_time.Stolen
	// 	fields[prefix+"cpu_time_guest"] = cpu_time.Guest
	// 	fields[prefix+"cpu_time_guest_nice"] = cpu_time.GuestNice
	// }

	cpu_perc, err := proc.Percent(time.Duration(1)*time.Second)
	if err == nil {
        fields[prefix+"cpu_usage"] = cpu_perc
	}

	mem, err := proc.MemoryInfo()
	if err == nil {
		fields[prefix+"memory_rss"] = mem.RSS
		fields[prefix+"memory_vms"] = mem.VMS
		fields[prefix+"memory_swap"] = mem.Swap
		fields[prefix+"memory_data"] = mem.Data
		fields[prefix+"memory_stack"] = mem.Stack
		fields[prefix+"memory_locked"] = mem.Locked
	}

	// rlims, err := proc.RlimitUsage(true)
	// if err == nil {
	// 	for _, rlim := range rlims {
	// 		var name string
	// 		switch rlim.Resource {
	// 		case process.RLIMIT_CPU:
	// 			name = "cpu_time"
	// 		case process.RLIMIT_DATA:
	// 			name = "memory_data"
	// 		case process.RLIMIT_STACK:
	// 			name = "memory_stack"
	// 		case process.RLIMIT_RSS:
	// 			name = "memory_rss"
	// 		case process.RLIMIT_NOFILE:
	// 			name = "num_fds"
	// 		case process.RLIMIT_MEMLOCK:
	// 			name = "memory_locked"
	// 		case process.RLIMIT_AS:
	// 			name = "memory_vms"
	// 		case process.RLIMIT_LOCKS:
	// 			name = "file_locks"
	// 		case process.RLIMIT_SIGPENDING:
	// 			name = "signals_pending"
	// 		case process.RLIMIT_NICE:
	// 			name = "nice_priority"
	// 		case process.RLIMIT_RTPRIO:
	// 			name = "realtime_priority"
	// 		default:
	// 			continue
	// 		}

	// 		fields[prefix+"rlimit_"+name+"_soft"] = rlim.Soft
	// 		fields[prefix+"rlimit_"+name+"_hard"] = rlim.Hard
	// 		if name != "file_locks" { // gopsutil doesn't currently track the used file locks count
	// 			fields[prefix+name] = rlim.Used
	// 		}
	// 	}
    // }

	// fmt.Println(fields)

	fmt.Printf("service_cpu,zone=%d,service=%d usage=%v\n", 
				zoneId, service, cpu_perc)

	if mem != nil {
		fmt.Printf("service_mem,zone=%d,service=%d rss=%v,vms=%v,swap=%v,data=%v,stack=%v,locked=%v\n",
					zoneId, service, 
					mem.RSS, mem.VMS, mem.Swap, mem.Data, mem.Stack, mem.Locked)
	}
}

func main ()  {
	if len(os.Args) > 1 {
		service_dir = os.Args[1]
	}

	// fmt.Printf("Path:%v\n", service_dir)

	loadServiceInfo()

	// fmt.Println(zoneId)
	for service, pid := range services {
		// fmt.Println(service, pid)
		go printServiceMetrics(service, pid)
	}

	for range services {
		<- done
	}
}