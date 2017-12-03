/*
Author: lkong
Description: test cmd tool
*/

package cmd

import (
	"encoding/json"
	"io"
	"strconv"

	"grapehttp/client/cmd/templates"
	cmdutil "grapehttp/client/cmd/util"
	"grapehttp/pkg/i18n"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	listExample = templates.Examples(i18n.T(`
	# List existing users with default mysql limit = 500 and offset = 0
	fctl list

	# List existing users with mysql limit = 10
	fctl list 10

	# List existing users with mysql limit = 10 and offset = 2
	fctl list 10 2`))
)

func NewCmdUserList(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   i18n.T("List existing users with mysql limit and offset"),
		Long:    "List existing users with mysql limit and offset",
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunUserList(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{"li"},
	}

	return cmd
}

func RunUserList(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	req := struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}{0, 500}

	if len(args) >= 1 {
		req.Limit, _ = strconv.Atoi(args[0])
	}

	if len(args) >= 2 {
		req.Offset, _ = strconv.Atoi(args[1])
	}

	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/list").
		Set("Authorization", "Basic "+f.Auth()).
		Send(req).
		End()

	if err := cmdutil.CombineRequestErr(resp, body, errs); err != nil {
		return err
	}

	users := []user{}

	if err := json.Unmarshal([]byte(body), &users); err != nil {
		return err
	}

	table := tablewriter.NewWriter(out)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(TABLE_WIDTH)
	table.SetHeader([]string{"Username", "Nickname", "Email", "Status", "Logincount", "Lastlogintime"})
	for _, user := range users {
		status := strconv.Itoa(user.Status)
		if user.Status == 0 {
			status = color.RedString("0")
		}

		table.Append([]string{user.Username, user.Nickname, user.Email, status,
			strconv.Itoa(user.Logincount), user.Lastlogintime.Format("2006-01-02 15:04:05")})
	}
	table.Render()
	return nil
}
