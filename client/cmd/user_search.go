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
	searchExample = templates.Examples(i18n.T(`
		# Search users which contains string: lk
		fctl search lk`))
)

func NewCmdUserSearch(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search CONDITION",
		Short:   i18n.T("Search users with fuzzy match condition"),
		Long:    "Search users with fuzzy match condition",
		Example: searchExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunUserSearch(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{"se"},
	}

	return cmd
}

func RunUserSearch(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	req := struct {
		Username string `json:"username"`
	}{args[0]}
	request := f.Gorequest()
	resp, body, errs := request.Get("http://"+f.Server+"/-/user/search").
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
