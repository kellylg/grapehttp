/*
Author: lkong
Description: test cmd tool
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"grapehttp/client/cmd/templates"
	cmdutil "grapehttp/client/cmd/util"
	"grapehttp/pkg/i18n"

	"github.com/spf13/cobra"
)

type HTTPStaticServer struct {
	Root            string
	Upload          bool
	Delete          bool
	Title           string
	Theme           string
	PlistProxy      string
	GoogleTrackerId string
	AuthType        string
}

var (
	lsExample = templates.Examples(i18n.T(`
	# List directory contents
	fctl ls /
	
	# More usage, type: fctl ls -h`))
)

func NewCmdLs(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls PATH",
		Short:   "List directory contents",
		Long:    "List directory contents",
		Example: lsExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				args = append(args, "/")
			}
			cmdutil.CheckErr(RunLs(f, out, cmdErr, cmd, args))
			return
		},
	}

	cmd.Flags().BoolP("all", "a", false, "do not ignore entries starting with .")
	cmd.Flags().BoolP("almost-all", "A", false, "do not list implied . and ..")
	cmd.Flags().BoolP("ignore-backups", "B", false, "do not list implied entries ending with ~")
	cmd.Flags().BoolP("column", "C", true, "list entries by columns")
	cmd.Flags().BoolP("color", "", true, "colorize the output")
	cmd.Flags().BoolP("human-readable", "H", false, "with -l, print sizes in human readable format (e.g., 1K 234M 2G)")
	cmd.Flags().BoolP("reverse", "r", false, "reverse order while sorting")
	cmd.Flags().BoolP("recursive", "R", false, "list subdirectories recursively")
	cmd.Flags().BoolP("classify", "F", false, "append indicator (one of */=>@|) to entries")

	cmd.Flags().BoolP("size", "S", false, "sort by file size")
	cmd.Flags().BoolP("time", "t", false, "sort by modification time, newest first")
	cmd.Flags().BoolP("line", "1", false, "list one file per line")
	cmd.Flags().BoolP("long", "l", false, "use a long listing format")
	return cmd
}

func RunLs(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	cmdArgs := buildArgs(cmd)
	req := struct {
		Name  string   `json:"name"`
		Args  []string `json:"args"`
		Paths []string `json:"paths"`
	}{"ls", cmdArgs, args}

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

func httpServerInfo(f cmdutil.Factory) (*HTTPStaticServer, error) {
	request := f.Gorequest()
	s := &HTTPStaticServer{}
	resp, body, errs := request.Get("http://"+f.Server+"/-/status").
		Send(``).
		Set("Authorization", "Basic "+f.Auth()).
		EndBytes()

	if err := cmdutil.CombineRequestErr(resp, string(body), errs); err != nil {
		return nil, err
	}

	err := json.Unmarshal(body, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func buildArgs(cmd *cobra.Command) []string {
	cmdArgs := []string{}
	if cmdutil.GetFlagBool(cmd, "all") {
		cmdArgs = append(cmdArgs, "-a")
	}
	if cmdutil.GetFlagBool(cmd, "almost-all") {
		cmdArgs = append(cmdArgs, "-A")
	}
	if cmdutil.GetFlagBool(cmd, "ignore-backups") {
		cmdArgs = append(cmdArgs, "-B")
	}
	if cmdutil.GetFlagBool(cmd, "column") {
		cmdArgs = append(cmdArgs, "-C")
	}
	if cmdutil.GetFlagBool(cmd, "color") {
		cmdArgs = append(cmdArgs, "--color")
	}
	if cmdutil.GetFlagBool(cmd, "human-readable") {
		cmdArgs = append(cmdArgs, "-h")
	}
	if cmdutil.GetFlagBool(cmd, "reverse") {
		cmdArgs = append(cmdArgs, "-r")
	}
	if cmdutil.GetFlagBool(cmd, "recursive") {
		cmdArgs = append(cmdArgs, "-R")
	}
	if cmdutil.GetFlagBool(cmd, "classify") {
		cmdArgs = append(cmdArgs, "-F")
	}
	if cmdutil.GetFlagBool(cmd, "size") {
		cmdArgs = append(cmdArgs, "-S")
	}
	if cmdutil.GetFlagBool(cmd, "time") {
		cmdArgs = append(cmdArgs, "-t")
	}
	if cmdutil.GetFlagBool(cmd, "line") {
		cmdArgs = append(cmdArgs, "-1")
	}
	if cmdutil.GetFlagBool(cmd, "long") {
		cmdArgs = append(cmdArgs, "-l")
	}

	return cmdArgs
}
