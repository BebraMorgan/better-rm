package flags

import (
	"brm/actions"
	"brm/localization"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

type Options struct {
	Verbose         bool
	ShowHelp        bool
	ShowVersion     bool
	IFlag           bool
	InteractiveOnce bool
	EmptyTrash      bool
}

func ParseFlags() Options {
	var opts Options

	pflag.BoolVarP(&opts.IFlag, "interactive-each", "i", false, localization.GetMessage("flag_interactive_i"))
	pflag.BoolVarP(&opts.InteractiveOnce, "interactive-once", "I", false, localization.GetMessage("flag_interactive_I"))
	pflag.BoolVarP(&opts.Verbose, "verbose", "v", false, localization.GetMessage("flag_verbose"))
	pflag.BoolVar(&opts.ShowHelp, "help", false, localization.GetMessage("flag_help"))
	pflag.BoolVar(&opts.ShowVersion, "version", false, localization.GetMessage("flag_version"))
	pflag.BoolVarP(&opts.EmptyTrash, "empty-trash", "e", false, localization.GetMessage("flag_empty_trash"))

	pflag.Usage = func() {
		msg := localization.GetMessage("usage_header", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, msg)
		pflag.PrintDefaults()
	}

	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			pflag.Usage()
			os.Exit(0)
		}
	}

	pflag.Parse()

	if opts.ShowVersion {
		printVersionAndExit()
	}
	if opts.EmptyTrash {
		if err := actions.EmptyTrash(); err != nil {
			fmt.Fprintf(os.Stderr, "Error emptying trash: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	return opts
}

func Args() []string {
	return pflag.Args()
}

func printVersionAndExit() {
	const version = "brm 1.0.0"
	fmt.Println(version)
	os.Exit(0)
}

