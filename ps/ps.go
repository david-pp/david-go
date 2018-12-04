package main

import (
    "fmt"
    "time"
    "os"
    "strconv"
    // "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/process"
)

// Add metrics a single Process
func printMetrics(proc process.Process) {
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

	cpu_time, err := proc.Times()
	if err == nil {
		fields[prefix+"cpu_time_user"] = cpu_time.User
		fields[prefix+"cpu_time_system"] = cpu_time.System
		fields[prefix+"cpu_time_idle"] = cpu_time.Idle
		fields[prefix+"cpu_time_nice"] = cpu_time.Nice
		fields[prefix+"cpu_time_iowait"] = cpu_time.Iowait
		fields[prefix+"cpu_time_irq"] = cpu_time.Irq
		fields[prefix+"cpu_time_soft_irq"] = cpu_time.Softirq
		fields[prefix+"cpu_time_steal"] = cpu_time.Steal
		fields[prefix+"cpu_time_stolen"] = cpu_time.Stolen
		fields[prefix+"cpu_time_guest"] = cpu_time.Guest
		fields[prefix+"cpu_time_guest_nice"] = cpu_time.GuestNice
	}

	cpu_perc, err := proc.Percent(time.Duration(0))
	if err == nil {
        fields[prefix+"cpu_usage"] = cpu_perc
        // fmt.Println("....")
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

    p,err := proc.CPUPercent()
    

    // fmt.Println(fields)
    fmt.Printf("server,id=21 cpu=%v,mem=%v",p , fields[prefix+"memory_rss"])
}


func main() {
    // v, _ := mem.VirtualMemory()

    // // // almost every return value is a struct
    // fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

    // // // convert to JSON. String() is also implemented
	// fmt.Println(v)
	
	for {
		
	}

    if len(os.Args) < 2 {
        fmt.Println("Usage: exe pid")
        return
    }

    // fmt.Println(os.Args)//打印切片内容
    // for i := 0; i < len(os.Args); i++ {
    //     fmt.Println(os.Args[i])
    // }

    pid, err := strconv.ParseInt(os.Args[1], 10, 32)
    ret, err := process.NewProcess(int32(pid))
	if err != nil {
        fmt.Println(err,ret)
		// t.Errorf("error %v", err)
	} else {
        printMetrics(*ret)
    }
}