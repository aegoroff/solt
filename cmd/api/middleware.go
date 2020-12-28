package api

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/dustin/go-humanize"
	"log"
	"runtime"
	"runtime/pprof"
	"time"
)

type executorMemUsage struct {
	wrapped Executor
	c       *Conf
}

type executorTimeMeasure struct {
	wrapped Executor
	c       *Conf
	start   time.Time
}

type executorCPUProfile struct {
	wrapped Executor
	c       *Conf
}

type executorMemoryProfile struct {
	wrapped Executor
	c       *Conf
}

func newMemUsageExecutor(e Executor, c *Conf) Executor {
	em := executorMemUsage{
		wrapped: e,
		c:       c,
	}
	return &em
}

func newTimeMeasureExecutor(e Executor, c *Conf) Executor {
	em := executorTimeMeasure{
		wrapped: e,
		c:       c,
		start:   time.Now(),
	}
	return &em
}

func newCPUProfileExecutor(e Executor, c *Conf) Executor {
	em := executorCPUProfile{
		wrapped: e,
		c:       c,
	}
	return &em
}

func newMemoryProfileExecutor(e Executor, c *Conf) Executor {
	em := executorMemoryProfile{
		wrapped: e,
		c:       c,
	}
	return &em
}

func (e *executorMemUsage) Execute() error {
	err := e.wrapped.Execute()

	if !*e.c.Diag() {
		return err
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	e.c.Prn().Cprint("\n<gray>Alloc =</> <green>%s</>", humanize.IBytes(m.Alloc))
	e.c.Prn().Cprint("\t<gray>TotalAlloc =</> <green>%s</>", humanize.IBytes(m.TotalAlloc))
	e.c.Prn().Cprint("\t<gray>Sys =</> <green>%s</>", humanize.IBytes(m.Sys))
	e.c.Prn().Cprint("\t<gray>NumGC =</> <green>%v</>\n", m.NumGC)
	return err
}

func (e *executorTimeMeasure) Execute() error {
	err := e.wrapped.Execute()

	if !*e.c.Diag() {
		return err
	}

	elapsed := time.Since(e.start)
	e.c.Prn().Cprint("<gray>Working time:</> <green>%v</>\n", elapsed)
	return err
}

func (e *executorCPUProfile) Execute() error {
	if *e.c.Diag() && *e.c.Cpu() != "" {
		f, err := e.c.Fs().Create(*e.c.Cpu())
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
	if *e.c.Diag() && *e.c.Memory() != "" {
		f, err := e.c.Fs().Create(*e.c.Memory())
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
