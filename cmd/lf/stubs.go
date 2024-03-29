package lf

import "solt/internal/out"

type nullRemover struct{}

func (*nullRemover) Remove([]string) {}

type nullExister struct{}

func (*nullExister) MissingCount() int64               { return 0 }
func (*nullExister) Print(out.Printer, string, string) {}
func (*nullExister) Validate(string, []string)         {}
