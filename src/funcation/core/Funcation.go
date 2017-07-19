package core

import (
	"fmt"
	"funcation/golog"
	"os"
	"time"
)

func InitCpu() {
	for {
		UpdateCpuStat()
		time.Sleep(1 * time.Second)
	}
}

var cpuchain chan *CpuInfo = make(chan *CpuInfo)
var memchain chan *MemInfoMap = make(chan *MemInfoMap)

func Collect() {
	t := time.NewTicker(60 * time.Second).C

	for {
		<-t

		c := CpuMetrics()
		cpuchain <- c
		b := MemMetrics()
		memchain <- b
	}
}

func DealChain() {
	cpufilefd, err := golog.NewFileHandler(CPUFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		golog.Error("core", "core", "OPEN CPUfile faild!", 0)
		os.Exit(1)
	}
	memfilefd, err := golog.NewFileHandler(MEMFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		golog.Error("core", "core", "OPEN MEMfile faild!", 0)
		os.Exit(1)
	}

	go func() {
		for {
			a := <-cpuchain
			cpufilefd.Write([]byte(fmt.Sprintf("%d ", a.Clicktime)))
			for _, v := range a.cpusidle {
				cpufilefd.Write([]byte(fmt.Sprintf("%s ", v)))
			}
			cpufilefd.Write([]byte(fmt.Sprintf("\n")))
		}
		defer cpufilefd.Close()
	}()
	go func() {
		for {
			b := <-memchain
			memfilefd.Write([]byte(fmt.Sprintln(b.Clicktime, b.pmemfree, b.pmemused, b.pswapfree, b.pswapused)))
		}
		defer memfilefd.Close()
	}()
}
