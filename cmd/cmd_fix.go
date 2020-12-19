package cmd

import "github.com/spf13/cobra"

type fixCommand struct {
	baseCommand
}

type sdkProjectReference struct {
	Path string `xml:"Include,attr"`
}

func newFix(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() command {
			return &fixCommand{
				baseCommand: newBaseCmd(c),
			}
		},
	}

	cmd := cc.newCobraCommand("fix", "fixprojects", "Fixes redundant SDK projects references")

	return cmd
}

func (c *fixCommand) execute() error {
	v := newsdkProjectsFixer(c.prn, c.fs)
	m := newSdkProjectsModule(c.fs, c.prn, c.sourcesPath, v)

	m.execute()
	return nil
}
