package cli

import (
	"os"
	"runtime/debug"

	"github.com/spf13/pflag"
)

var (
	help      bool = false
	dryRun    bool = false
	debugMode bool = false
	workDir   string
)

func init() {
	pwd, _ := os.Getwd()
	pflag.StringVarP(&workDir, "workdir", "w", pwd, "Working directory")
	pflag.BoolVarP(&help, "help", "h", false, "Print usage")
	pflag.BoolVarP(&dryRun, "dry-run", "n", false, "Perform a trial run with no changes made")
	pflag.BoolVarP(&debugMode, "debug", "", false, "Enable debug logging")
}

// TODO pass in options?
func Recast(fn func(func())) {
	pflag.Parse()
	if help {
		pflag.Usage()
		os.Exit(0)
	}
	if debugMode {
		log.SetLevel(log.DEBUG)
	}
	core.DryRun = dryRun
	core.WorkingDir = workDir

	defer func() {
		if err := recover(); err != nil {
			if errors.IsRecastError(err) {
				log.Errorf("Error: %v", err)
			} else {
				log.Errorf("Unexpected Error: %v", err)
			}
			// TODO configure stack printing
			log.Debugf(string(debug.Stack()))
		}
	}()

	fn(ExecRegistred)
}
