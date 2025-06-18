package main

import (
	"brm/actions"
	"brm/flags"
	"brm/localization"
	"brm/tui/browser"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strings"
)

func confirmPrompt(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	answer := strings.ToLower(strings.TrimSpace(result))
	return answer == "y" || answer == "", nil
}

func deleteWithConfirmation(args []string, opts flags.Options) {
	if opts.InteractiveOnce && len(args) > 3 {
		confirmed, err := confirmPrompt(localization.GetMessage("confirm_delete_files", len(args)))
		if err != nil {
			fmt.Println(localization.GetMessage("delete_cancelled"))
			return
		}
		if !confirmed {
			fmt.Println(localization.GetMessage("delete_cancelled"))
			return
		}
		for _, arg := range args {
			deleteFile(arg, opts.Verbose)
		}
	} else {
		for _, arg := range args {
			if opts.IFlag {
				confirmed, err := confirmPrompt(localization.GetMessage("confirm_delete_file", arg))
				if err != nil {
					fmt.Println(localization.GetMessage("delete_cancelled"))
					return
				}
				if !confirmed {
					fmt.Println(localization.GetMessage("delete_cancelled"))
					return
				}
			}
			deleteFile(arg, opts.Verbose)
		}
	}
}

func deleteFile(arg string, verbose bool) {
	err := actions.SaveDelete(arg)
	if err != nil {
		log.Printf(localization.GetMessage("error_moving_to_trash"), arg, err)
	} else {
		if verbose {
			fmt.Println(localization.GetMessage("file_deleted_verbose", arg))
		}
	}
}

func main() {
	opts := flags.ParseFlags()
	args := flags.Args()

	if len(args) == 0 && !opts.ShowHelp && !opts.ShowVersion && !opts.EmptyTrash {

		p := tea.NewProgram(browser.NewModel(""))
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
			os.Exit(1)
		}

		return
	}

	deleteWithConfirmation(args, opts)
}
