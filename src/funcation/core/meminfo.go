package core

import (
	"log"
	"time"

	"github.com/toolkits/nux"
)

type MemInfoMap struct {
	Clicktime int64
	Mem       *nux.Mem
	pmemfree  float64
	pmemused  float64
	pswapfree float64
	pswapused float64
}

func MemMetrics() *MemInfoMap {

	var l MemInfoMap
	m, err := nux.MemInfo()
	if err != nil {
		log.Println(err)
		return nil
	}

	memFree := m.MemFree + m.Buffers + m.Cached
	memUsed := m.MemTotal - memFree

	pmemFree := 0.0
	pmemUsed := 0.0
	if m.MemTotal != 0 {
		pmemFree = float64(memFree) * 100.0 / float64(m.MemTotal)
		pmemUsed = float64(memUsed) * 100.0 / float64(m.MemTotal)
	}

	pswapFree := 0.0
	pswapUsed := 0.0
	if m.SwapTotal != 0 {
		pswapFree = float64(m.SwapFree) * 100.0 / float64(m.SwapTotal)
		pswapUsed = float64(m.SwapUsed) * 100.0 / float64(m.SwapTotal)
	}
	l.Clicktime = time.Now().Unix()
	l.Mem = m
	l.pmemfree = pmemFree
	l.pmemused = pmemUsed
	l.pswapfree = pswapFree
	l.pswapused = pswapUsed
	return &l

}
