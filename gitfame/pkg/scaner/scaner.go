package scaner

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Scaner struct {
	Repository   string
	Revision     string
	OrderBy      string
	UseCommitter bool
	Format       string
	Extensions   string
	Languages    string
	Exclude      string
	RestrictTo   string
}

var Log *logrus.Logger

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("repository", "r", ".", "Path to Git repository")
	cmd.Flags().StringP("revision", "", "HEAD", "Git revision")
	cmd.Flags().StringP("order-by", "", "lines", "Sort results by 'lines', 'commits', or 'files'")
	cmd.Flags().BoolP("use-committer", "", false, "Use committer instead of author in calculations")
	cmd.Flags().StringP("format", "", "tabular", "Output format: 'tabular', 'csv', 'json', 'json-lines'")
	cmd.Flags().StringP("extensions", "", "", "List of file extensions to include")
	cmd.Flags().StringP("languages", "", "", "List of programming languages to include")
	cmd.Flags().StringP("exclude", "", "", "Glob patterns to exclude files")
	cmd.Flags().StringP("restrict-to", "", "", "Glob patterns to include files")
}

func readFlags(cmd *cobra.Command, s *Scaner) {
	s.Repository, _ = cmd.Flags().GetString("repository")
	s.Revision, _ = cmd.Flags().GetString("revision")
	s.OrderBy, _ = cmd.Flags().GetString("order-by")
	s.UseCommitter, _ = cmd.Flags().GetBool("use-committer")
	s.Format, _ = cmd.Flags().GetString("format")
	s.Extensions, _ = cmd.Flags().GetString("extensions")
	s.Languages, _ = cmd.Flags().GetString("languages")
	s.Exclude, _ = cmd.Flags().GetString("exclude")
	s.RestrictTo, _ = cmd.Flags().GetString("restrict-to")
}

func (s *Scaner) Scan(args []string) {
	Log = logrus.New()
	Log.SetLevel(logrus.DebugLevel)
	var rootCmd = &cobra.Command{
		Use: "gitfame",
		Run: func(cmd *cobra.Command, args []string) {
			readFlags(cmd, s)
		},
	}

	setFlags(rootCmd)

	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		Log.Fatal(err)
	}
	//Log.Debug("repository: ", s.Repository)
	//Log.Debug("revision: ", s.Revision)
	//Log.Debug("orderby: ", s.OrderBy)
	//Log.Debug("use-commiter: ", s.UseCommitter)
	//Log.Debug("format: ", s.Format)
	//Log.Debug("extensions: ", s.Extensions)
	//Log.Debug("languages:  ", s.Languages)
	//Log.Debug("exclude: ", s.Exclude)
	//Log.Debug("restrictTo: ", s.RestrictTo)
}
