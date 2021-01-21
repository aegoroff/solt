package nuget

import (
	"solt/internal/ux"
)

type totalsBySolution struct {
	solutions  int64
	nugets     int64
	mismatched int64
}

type totalsByProjects struct {
	projects int64
	nugets   int64
}

func (t *totalsBySolution) Display(tbl *ux.Tabler) {
	tbl.AddLines(ux.NewLine("Solutions", t.solutions))
	tbl.AddLine("", "")
	tbl.AddHead("Packages", "Count")

	tbl.AddLines(
		ux.NewLine("Total", t.nugets),
		ux.NewLine("Mismatched", t.mismatched),
	)
}

func (t *totalsByProjects) Display(tbl *ux.Tabler) {
	tbl.AddLines(
		ux.NewLine("Projects", t.projects),
		ux.NewLine("Packages", t.nugets),
	)
}
