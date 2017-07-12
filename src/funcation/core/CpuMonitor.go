package core

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/toolkits/nux"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

type CpuInfo struct {
	Clicktime int64
	cpusidle  []string
}

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

//获取cpu - processsor  idle time 总数
func idleTotal() []uint64 {

	if procStatHistory[1] == nil {
		return make([]uint64, len(procStatHistory[0].Cpus))
	}
	psc_total := make([]uint64, len(procStatHistory[0].Cpus))
	for i, _ := range procStatHistory[1].Cpus {
		psc_total[i] = procStatHistory[0].Cpus[i].Total - procStatHistory[1].Cpus[i].Total
		//psc_total = append(psc_total, procStatHistory[0].Cpus[i].Total-procStatHistory[1].Cpus[i].Total)
	}
	return psc_total
}

func PsCpuidle(w []uint64) []string {
	var invQuotient float64
	if procStatHistory[1] == nil || procStatHistory[0] == nil {
		return nil
	}
	psLock.RLock()
	defer psLock.RUnlock()
	psidle := make([]string, len(procStatHistory[0].Cpus))
	for i, v := range w {

		invQuotient = 100.00 / float64(v)
		psidle[i] = strconv.FormatFloat(float64((procStatHistory[0].Cpus[i].Idle-procStatHistory[1].Cpus[i].Idle))*invQuotient, 'f', 2, 64)
		//		psidle = append(psidle, float64((procStatHistory[0].Cpus[i].Idle-procStatHistory[1].Cpus[i].Idle))*invQuotient)
	}
	return psidle
}

func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle) * invQuotient
}

func CpuUser() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
}

func CpuNice() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice) * invQuotient
}

func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
}

func CpuIowait() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
}

func CpuIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq) * invQuotient
}

func CpuSoftIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq) * invQuotient
}

func CpuSteal() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal) * invQuotient
}

func CpuGuest() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest) * invQuotient
}

func CurrentCpuSwitches() uint64 {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[0].Ctxt
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

/*func CpuMetrics() []*model.MetricValue {
	if !CpuPrepared() {
		return []*model.MetricValue{}
	}

	cpuIdleVal := CpuIdle()
	idle := GaugeValue("cpu.idle", cpuIdleVal)
	busy := GaugeValue("cpu.busy", 100.0-cpuIdleVal)
	user := GaugeValue("cpu.user", CpuUser())
	nice := GaugeValue("cpu.nice", CpuNice())
	system := GaugeValue("cpu.system", CpuSystem())
	iowait := GaugeValue("cpu.iowait", CpuIowait())
	irq := GaugeValue("cpu.irq", CpuIrq())
	softirq := GaugeValue("cpu.softirq", CpuSoftIrq())
	steal := GaugeValue("cpu.steal", CpuSteal())
	guest := GaugeValue("cpu.guest", CpuGuest())
	switches := CounterValue("cpu.switches", CurrentCpuSwitches())
	return []*model.MetricValue{idle, busy, user, nice, system, iowait, irq, softirq, steal, guest, switches}
}*/

func CpuMetrics() *CpuInfo {
	var l CpuInfo
	if !CpuPrepared() {
		fmt.Println("hello -world")
		return nil
	}
	t := idleTotal()
	l.Clicktime = time.Now().Unix()
	l.cpusidle = PsCpuidle(t)

	return &l
}

func Cputest() {
	for i := 10; i > 0; i-- {
		t := idleTotal()
		for index, v := range PsCpuidle(t) {
			fmt.Printf("cpu%d: %f\n", index, v)
		}
		fmt.Println("<>>>>>>>>>>>>>>>>>>>>><>")
	}
}
