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
	"grapehttp/pkg/i18n"

	"github.com/spf13/cobra"
)

var (
	cpExample = templates.Examples(i18n.T(`
		# Copy files and directories
		gctl cp /A /

		# Do not prompt before overwriting.
		gctl cp -f /A /

		# Copy directories recursively
		gctl cp -r /dir1 /dir2`))
)

func NewCmdCp(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cp SRCFILE DSTPATH",
		Short:   "Copy files and directories",
		Long:    "Copy files and directories",
		Example: cpExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunCp(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().BoolP("force", "f", false, "do not prompt before overwriting.")
	cmd.Flags().BoolP("archive", "a", false, "same as -dR --preserve=all.")
	cmd.Flags().BoolP("recursive", "r", false, "copy directories recursively.")
	return cmd
}

func RunCp(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	cmdArgs := []string{}
	if cmdutil.GetFlagBool(cmd, "force") {
		cmdArgs = append(cmdArgs, "-f")
	}
	if cmdutil.GetFlagBool(cmd, "archive") {
		cmdArgs = append(cmdArgs, "-a")
	}
	if cmdutil.GetFlagBool(cmd, "recursive") {
		cmdArgs = append(cmdArgs, "-r")
	}

	req := struct {
		Name  string   `json:"name"`
		Args  []string `json:"args"`
		Paths []string `json:"paths"`
	}{"cp", cmdArgs, args}

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
