package cmd

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/dustin/go-humanize"
	"log"
	"runtime"
	"runtime/pprof"
	"solt/cmd/api"
	"time"
)

type executorMemUsage struct {
	wrapped api.Executor
	c       *conf
}

type executorTimeMeasure struct {
	wrapped api.Executor
	c       *conf
	start   time.Time
}

type executorCPUProfile struct {
	wrapped api.Executor
	c       *conf
}

type executorMemoryProfile struct {
	wrapped api.Executor
	c       *conf
}

func newMemUsageExecutor(e api.Executor, c *conf) api.Executor {
	em := executorMemUsage{
		wrapped: e,
		c:       c,
	}
	return &em
}

func newTimeMeasureExecutor(e api.Executor, c *conf) api.Executor {
	em := executorTimeMeasure{
		wrapped: e,
		c:       c,
		start:   time.Now(),
	}
	return &em
}

func newCPUProfileExecutor(e api.Executor, c *conf) api.Executor {
	em := executorCPUProfile{
		wrapped: e,
		c:       c,
	}
	return &em
}

func newMemoryProfileExecutor(e api.Executor, c *conf) api.Executor {
	em := executorMemoryProfile{
		wrapped: e,
		c:       c,
	}
	return &em
}

func (e *executorMemUsage) Execute() error {
	err := e.wrapped.Execute()

	if !*e.c.diag {
		return err
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	e.c.p.Cprint("\n<gray>Alloc =</> <green>%s</>", humanize.IBytes(m.Alloc))
	e.c.p.Cprint("\t<gray>TotalAlloc =</> <green>%s</>", humanize.IBytes(m.TotalAlloc))
	e.c.p.Cprint("\t<gray>Sys =</> <green>%s</>", humanize.IBytes(m.Sys))
	e.c.p.Cprint("\t<gray>NumGC =</> <green>%v</>\n", m.NumGC)
	return err
}

func (e *executorTimeMeasure) Execute() error {
	err := e.wrapped.Execute()

	if !*e.c.diag {
		return err
	}

	elapsed := time.Since(e.start)
	e.c.p.Cprint("<gray>Working time:</> <green>%v</>\n", elapsed)
	return err
}

func (e *executorCPUProfile) Execute() error {
	if *e.c.diag && *e.c.cpu != "" {
		f, err := e.c.fs().Create(*e.c.cpu)
		if err != nil {
			return err
		}
		err = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	return e.wrapped.Execute()
}

func (e *executorMemoryProfile) Execute() error {
	err := e.wrapped.Execute()
	if *e.c.diag && *e.c.memory != "" {
		f, err := e.c.fs().Create(*e.c.memory)
		if err != nil {
			return err
		}
		defer scan.Close(f)

		err = pprof.WriteHeapProfile(f)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}
