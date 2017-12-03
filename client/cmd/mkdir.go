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
	mkdirExample = templates.Examples(i18n.T(`
	# Make directories
	fctl mkdir /lkong`))
)

func NewCmdMkdir(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mkdir PATH",
		Short:   "Create specified directory",
		Long:    "Create specified directory",
		Example: mkdirExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunMkdir(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{"mk"},
	}

	cmd.Flags().BoolP("parents", "p", false, "no error if existing, make parent directories as needed.")
	return cmd
}

func RunMkdir(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	cmdArgs := []string{}
	if cmdutil.GetFlagBool(cmd, "parents") {
		cmdArgs = append(cmdArgs, "-p")
	}

	req := struct {
		Name  string   `json:"name"`
		Args  []string `json:"args"`
		Paths []string `json:"paths"`
	}{"mkdir", cmdArgs, args}

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
