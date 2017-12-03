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
	getExample = templates.Examples(i18n.T(`
		# Get lkong user informations
		fctl get lkong`))
)

func NewCmdUserGet(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get USERNAME",
		Short:   i18n.T("Get user informations"),
		Long:    "Get user informations",
		Example: getExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunUserGet(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	return cmd
}

func RunUserGet(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	req := struct {
		Username string `json:"username"`
	}{args[0]}
	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/get").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
