/*
Author: lkong
Description: test cmd tool
*/

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"grapehttp/client/cmd/templates"
	cmdutil "grapehttp/client/cmd/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type DownloadOptions struct {
	output string
}

var (
	downloadExample = templates.Examples(`
		# Download test.log to current directory
		fctl download /api.log

		# Download test.log to specified directory
		fctl download /api.log -o /data

		# Download multiple files from different directory
		fctl download /test.txt /lkong/api.log`)
)

func NewCmdDownload(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download REMOTE_FILE LOCAL_FILE",
		Short:   "Download files from remote http server",
		Long:    "Download files from remote http server",
		Example: downloadExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			options := new(DownloadOptions)
			cmdutil.CheckErr(options.Complete(cmd))
			cmdutil.CheckErr(options.Run(f, out, cmdErr, args))
			return
		},
		Aliases: []string{"do"},
	}

	cmd.Flags().StringP("output", "o", ".", "download file to `output` directory, default .")
	return cmd
}

func (o *DownloadOptions) Complete(cmd *cobra.Command) error {
	o.output = cmdutil.GetFlagString(cmd, "output")
	return nil
}

func (o *DownloadOptions) Run(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, args []string) error {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))

	for _, name := range args {
		url := "http://" + f.Server + name
		wg.Add(1)
		go o.download(f, &wg, p, path.Base(name), url)
	}

	wg.Wait()
	p.Stop()
	return nil
}

func (o DownloadOptions) download(f cmdutil.Factory, wg *sync.WaitGroup, p *mpb.Progress, name, url string) error {
	defer wg.Done()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Yellow("%s: %v", name, err)
		return err
	}

	req.Header.Set("Authorization", "Basic "+f.Auth())
	req.Header.Set("Username", f.Username)
	resp, err := client.Do(req)
	if err != nil {
		color.Yellow("%s: %v", name, err)
		return fmt.Errorf("%s: %v", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Yellow("%s: %v", name, resp.Status)
		return fmt.Errorf("%s: %v", name, err)
	}

	size := resp.ContentLength

	// create dest
	destName := filepath.Base(url)
	dest, err := os.Create(filepath.Join(o.output, destName))
	if err != nil {
		err = fmt.Errorf("Can't create %s: %v", destName, err)
		color.Yellow("%s: %v", name, resp.Status)
		return fmt.Errorf("%s: %v", name, err)
	}

	// create bar with appropriate decorators
	bar := p.AddBar(size,
		mpb.PrependDecorators(
			decor.StaticName(color.CyanString(name), 0, 0),
			decor.Counters("%3s / %3s", decor.Unit_KiB, 18, decor.DSyncSpace),
		),
		mpb.AppendDecorators(decor.Percentage(5, 0)))

	if !f.Cool {
		p.RemoveBar(bar)
	}

	// create proxy reader
	reader := bar.ProxyReader(resp.Body)
	// and copy from reader
	_, err = io.Copy(dest, reader)

	if closeErr := dest.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		p.RemoveBar(bar)
		color.Yellow("%s: %v", name, err)
		return fmt.Errorf("%s: %v", name, err)
	}
	//p.RemoveBar(bar)
	return nil
}
