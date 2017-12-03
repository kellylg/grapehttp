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
	disableExample = templates.Examples(i18n.T(`
		# Disable users
		fctl disable user1 user2`))
)

func NewCmdUserDisable(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable USERNAME [USERNAME]",
		Short:   i18n.T("Disable users"),
		Long:    "Disable users",
		Example: disableExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunUserDisable(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	return cmd
}

func RunUserDisable(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	req := struct {
		Usernames []string `json:"usernames"`
	}{args}
	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/disable").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s", body)
	return nil
}
