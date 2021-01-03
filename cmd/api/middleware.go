package api

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"log"
	"runtime"
	"runtime/pprof"
	"time"
)

type executorMiddleware struct {
	wrapped Executor
	c       *Conf
}

func newExecutorMiddleware(wrapped Executor, c *Conf) *executorMiddleware {
	return &executorMiddleware{wrapped: wrapped, c: c}
}

type executorMemUsage struct {
	*executorMiddleware
}

type executorTimeMeasure struct {
	*executorMiddleware
	start time.Time
}

type executorCPUProfile struct {
	*executorMiddleware
}

type executorMemoryProfile struct {
	*executorMiddleware
}

func newMemUsageExecutor(e Executor, c *Conf) Executor {
	em := executorMemUsage{
		newExecutorMiddleware(e, c),
	}
	return &em
}

func newTimeMeasureExecutor(e Executor, c *Conf) Executor {
	em := executorTimeMeasure{
		executorMiddleware: newExecutorMiddleware(e, c),
		start:              time.Now(),
	}
	return &em
}

func newCPUProfileExecutor(e Executor, c *Conf) Executor {
	em := executorCPUProfile{
		newExecutorMiddleware(e, c),
	}
	return &em
}

func newMemoryProfileExecutor(e Executor, c *Conf) Executor {
	em := executorMemoryProfile{
		newExecutorMiddleware(e, c),
	}
	return &em
}

func (e *executorMemUsage) Execute(cc *cobra.Command) error {
	err := e.wrapped.Execute(cc)

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

func (e *executorTimeMeasure) Execute(cc *cobra.Command) error {
	err := e.wrapped.Execute(cc)

	if !*e.c.Diag() {
		return err
	}

	elapsed := time.Since(e.start)
	e.c.Prn().Cprint("<gray>Working time:</> <green>%v</>\n", elapsed)
	return err
}

func (e *executorCPUProfile) Execute(cc *cobra.Command) error {
	if *e.c.Diag() && *e.c.CPU() != "" {
		f, err := e.c.Fs().Create(*e.c.CPU())
		if err != nil {
			return err
		}
		err = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	return e.wrapped.Execute(cc)
}

func (e *executorMemoryProfile) Execute(cc *cobra.Command) error {
	err := e.wrapped.Execute(cc)
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
