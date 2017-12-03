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
	enableExample = templates.Examples(i18n.T(`
		# Enable users
		fctl enable user1 user2`))
)

func NewCmdUserEnable(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enable USERNAME [USERNAME]",
		Short:   i18n.T("Enable users"),
		Long:    "Enable users",
		Example: enableExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunUserEnable(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	return cmd
}

func RunUserEnable(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	req := struct {
		Usernames []string `json:"usernames"`
	}{args}
	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/enable").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s", body)
	return nil
}
