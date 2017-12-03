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

type UserModifyOptions struct {
	password string
	nickname string
	email    string
	remark   string
	status   int
}

var (
	modifyExample = templates.Examples(i18n.T(`
		# Modify user lkong's email address
		fctl modify lkong -e newaddress@tencent.com

		# Modify user lkong's password
		fctl modify lkong -p newpassword

		# Modify user lkong's nickname
		fctl modify lkong -n newnickname

		# Modify user lkong's remark
		fctl modify lkong -r newremark`))
)

func NewCmdUserModify(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "modify USERNAME",
		Short:   i18n.T("Modify user's email, password, nickname and remark"),
		Long:    "Modify user's email, password, nickname and remark",
		Example: modifyExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(validateUserModifyArgs(cmd, args))
			options := new(UserModifyOptions)
			cmdutil.CheckErr(options.Complete(cmd))
			if err := options.Validate(); err != nil {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, err.Error()))
			}
			cmdutil.CheckErr(options.RunUserModify(f, out, cmdErr, args))
			return
		},
		Aliases: []string{},
	}

	cmd.Flags().StringP("password", "p", "", "Specify the user password.")
	cmd.Flags().StringP("email", "e", "", "Specify the user email.")
	cmd.Flags().StringP("remark", "r", "", "Specify the user remark.")
	cmd.Flags().StringP("nickname", "n", "", "Specify the user nickname.")
	return cmd
}

func validateUserModifyArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if !isUsername(args[0]) {
		return fmt.Errorf("%s is not a legal user name", args[0])
	}

	return nil
}

func (o *UserModifyOptions) RunUserModify(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, args []string) error {
	req := user{
		Username: args[0],
		Password: o.password,
		Email:    o.email,
		Nickname: o.nickname,
		Remark:   o.remark,
	}

	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/modify").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	fmt.Printf("%s", body)
	return nil
}

func (o *UserModifyOptions) Complete(cmd *cobra.Command) error {
	o.password = cmdutil.GetFlagString(cmd, "password")
	o.nickname = cmdutil.GetFlagString(cmd, "nickname")
	o.email = cmdutil.GetFlagString(cmd, "email")
	o.remark = cmdutil.GetFlagString(cmd, "remark")
	return nil
}

func (o *UserModifyOptions) Validate() error {
	if o.email != "" {
		if !isEmail(o.email) {
			return fmt.Errorf("%s is not a email format", o.email)
		}
	}

	return nil
}
