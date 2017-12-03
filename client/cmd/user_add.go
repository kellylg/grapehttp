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
	addExample = templates.Examples(i18n.T(`
		# Add a new user lkong
		fctl add lkong pwd1234 lkong@tencent.com 

		# Add a new user lkong with nickname and remark
		fctl add lkong pwd1234 lkong@tencent.com -n lkong -r "test for add sub command"`))
)

func NewCmdUserAdd(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add USERNAME PASSWORD EMAIL",
		Short:   i18n.T("Create a new user"),
		Long:    "Create a new user",
		Example: addExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(validateUserArgs(cmd, args))
			cmdutil.CheckErr(RunUserAdd(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().StringP("nickname", "n", "", "Specify the user nickname.")
	cmd.Flags().StringP("remark", "r", "", "Specify the user remark.")
	return cmd
}

func validateUserArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if !isUsername(args[0]) {
		fmt.Printf("%v\n", args)
		return fmt.Errorf("%s is not a valid user name", args[0])
	}

	if !isEmail(args[2]) {
		return fmt.Errorf("%s is not a valid email name", args[2])
	}

	return nil
}

func RunUserAdd(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	nickname := cmdutil.GetFlagString(cmd, "nickname")
	if nickname == "" {
		nickname = args[0]
	}

	remark := cmdutil.GetFlagString(cmd, "remark")
	req := user{
		Username: args[0],
		Password: args[1],
		Email:    args[2],
		Nickname: nickname,
		Remark:   remark,
		Status:   1,
	}

	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/add").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s", body)
	return nil
}
