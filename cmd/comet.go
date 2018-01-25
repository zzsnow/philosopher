package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/comet"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// cometCmd represents the comet command
var cometCmd = &cobra.Command{
	Use:   "comet",
	Short: "Peptide spectrum matching with Comet",
	Run: func(cmd *cobra.Command, args []string) {

		// verify if the command is been executed on a workspace directory
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing Comet")

		m, e := comet.Run(m, args)
		if e != nil {
			logrus.Fatal(e.Error())
		}

		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "comet" {

		m.Restore(sys.Meta())

		cometCmd.Flags().BoolVarP(&m.Comet.Print, "print", "", false, "print a comet.params file")
		cometCmd.Flags().BoolVarP(&m.Comet.NoIndex, "noindex", "", false, "skip raw file indexing")
		cometCmd.Flags().StringVarP(&m.Comet.Param, "param", "", "comet.params.txt", "comet parameter file")
	}

	RootCmd.AddCommand(cometCmd)
}
