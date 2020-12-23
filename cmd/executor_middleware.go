package cmd

import (
	"github.com/dustin/go-humanize"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type executorMemUsage struct {
	wrapped executor
	c       *conf
}

type executorTimeMeasure struct {
	wrapped executor
	c       *conf
	start   time.Time
}

type executorCpuProfile struct {
	wrapped executor
	c       *conf
}

func newMemUsageExecutor(e executor, c *conf) executor {
	em := executorMemUsage{
		wrapped: e,
		c:       c,
	}
	return &em
}

func newTimeMeasureExecutor(e executor, c *conf) executor {
	em := executorTimeMeasure{
		wrapped: e,
		c:       c,
		start:   time.Now(),
	}
	return &em
}

func newCpuProfileExecutor(e executor, c *conf) executor {
	em := executorCpuProfile{
		wrapped: e,
		c:       c,
	}
	return &em
}

func (e *executorMemUsage) execute() error {
	err := e.wrapped.execute()

	if !*e.c.diag {
		return err
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	e.c.p.cprint("\n<gray>Alloc =</> <green>%s</>", humanize.IBytes(m.Alloc))
	e.c.p.cprint("\t<gray>TotalAlloc =</> <green>%s</>", humanize.IBytes(m.TotalAlloc))
	e.c.p.cprint("\t<gray>Sys =</> <green>%s</>", humanize.IBytes(m.Sys))
	e.c.p.cprint("\t<gray>NumGC =</> <green>%v</>\n", m.NumGC)
	return err
}

func (e *executorTimeMeasure) execute() error {
	err := e.wrapped.execute()

	if !*e.c.diag {
		return err
	}

	elapsed := time.Since(e.start)
	e.c.p.cprint("<gray>Working time:</> <green>%v</>\n", elapsed)
	return err
}

func (e *executorCpuProfile) execute() error {
	if *e.c.diag && *e.c.cpu != "" {
		f, err := os.Create(*e.c.cpu)
		if err != nil {
			return err
		}
		err = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	return e.wrapped.execute()
}
