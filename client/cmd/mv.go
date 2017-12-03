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
	mvExample = templates.Examples(i18n.T(`
		# Move (rename) files
		fctl mv /A /B`))
)

func NewCmdMv(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mv SRCPATH DSTPATH",
		Short:   "Move (rename) files",
		Long:    "Move (rename) files",
		Example: mvExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunMv(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().BoolP("force", "f", false, "do not prompt before overwriting.")
	return cmd
}

func RunMv(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	cmdArgs := []string{}
	if cmdutil.GetFlagBool(cmd, "force") {
		cmdArgs = append(cmdArgs, "-f")
	}

	req := struct {
		Name  string   `json:"name"`
		Args  []string `json:"args"`
		Paths []string `json:"paths"`
	}{"mv", cmdArgs, args}

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
