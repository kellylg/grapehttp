/*
Author: lkong
Description: test cmd tool
*/

package cmd

import (
	"fmt"
	"io"

	"grapehttp/client/cmd/templates"
	cmdutil "grapehttp/client/cmd/util"

	"github.com/spf13/cobra"
)

var (
	rmExample = templates.Examples(`
	# Command desc
	fctl rm

	# Command desc
	fctl rm -u lkong update

	# Command desc
	fctl rm -c download`)

	blackPath = []string{"", "/", "/bin", "/boot", "/dev", "/etc", "/home", "/lib", "/lib64", "/media", "/mnt", "/opt", "/proc", "/root", "/run", "/srv", "/sys", "/var", "/usr", "/data"}
)

func NewCmdRm(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm PATH [PATH]",
		Short:   "Remove files or directories",
		Long:    "Remove files or directories",
		Example: rmExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunRm(f, out, cmdErr, cmd, args))
			return
		},
	}

	cmd.Flags().BoolP("force", "f", false, "gnore nonexistent files and arguments, never prompt")
	cmd.Flags().BoolP("recursive", "r", false, "remove directories and their contents recursively")
	return cmd
}

func RunRm(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	cmdArgs := []string{}
	if cmdutil.GetFlagBool(cmd, "force") {
		cmdArgs = append(cmdArgs, "-f")
	}
	if cmdutil.GetFlagBool(cmd, "recursive") {
		cmdArgs = append(cmdArgs, "-r")
	}
	req := struct {
		Name  string   `json:"name"`
		Args  []string `json:"args"`
		Paths []string `json:"paths"`
	}{"rm", cmdArgs, args}

	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/cmd").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s", body)
	return nil
}
