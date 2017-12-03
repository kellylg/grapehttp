/*
Author: lkong
Description: test cmd tool
*/

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"grapehttp/client/cmd/templates"
	cmdutil "grapehttp/client/cmd/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var (
	uploadExample = templates.Examples(`
		# Upload multiple files to http server
		fctl upload test.txt api.txt /lkong`)
)

func NewCmdUpload(f cmdutil.Factory, out io.Writer, cmdErr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upload LOCAL_FILE [LOCAL_FILE] REMOTE_DIR",
		Short:   "Upload files to remote http server",
		Long:    "Upload files to remote http server",
		Example: uploadExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmdutil.CheckErr(cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args))
			}
			cmdutil.CheckErr(RunUpload(f, out, cmdErr, cmd, args))
			return
		},
		Aliases: []string{"up"},
	}

	return cmd
}

func RunUpload(f cmdutil.Factory, out io.Writer, cmdErr io.Writer, cmd *cobra.Command, args []string) error {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	dstDir := args[len(args)-1]
	url := "http://" + strings.Replace(f.Server+"/"+dstDir, "//", "/", -1)
	args = append(args[:len(args)-1], args[len(args):]...) //删除最后一个

	for _, file := range args {
		// make sure file is not dir
		pass, err := checkFile(file)
		if pass {
			wg.Add(1)
			//go upload(&wg, p, file, url)
			go upload(f, &wg, p, file, url)
		} else {
			color.Yellow("%v", err)
			continue
		}
	}

	wg.Wait()
	p.Stop()
	return nil
}

func upload(f cmdutil.Factory, wg *sync.WaitGroup, p *mpb.Progress, filename string, url string) error {
	name := filepath.Base(filename)
	defer wg.Done()
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", name)
	if err != nil {
		color.Yellow("%s: %v", name, err)
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	defer fh.Close()
	if err != nil {
		color.Yellow("%s: %v", name, err)
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		color.Yellow("%s: %v", name, err)
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	fileInfo, _ := fh.Stat()

	// create bar with appropriate decorators
	bar := p.AddBar(fileInfo.Size(),
		mpb.PrependDecorators(
			decor.StaticName(color.CyanString(name), 0, 0),
			decor.Counters("%3s / %3s", decor.Unit_KiB, 18, decor.DSyncSpace),
		),
		mpb.AppendDecorators(decor.Percentage(5, 0)),
	)

	if !f.Cool {
		p.RemoveBar(bar)
	}

	reader := bar.ProxyReader(bodyBuf)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		p.RemoveBar(bar)
		color.Yellow("%s: %v", name, err)
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Username", f.Username)
	req.Header.Set("Authorization", "Basic "+f.Auth())
	resp, err := client.Do(req)
	//resp, err := http.Post(url, contentType, reader)
	if err != nil {
		p.RemoveBar(bar)
		color.Yellow("%s: %v", name, err)
		return err
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.RemoveBar(bar)
		color.Yellow("%s: %v", name, err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		p.RemoveBar(bar)
		//color.Yellow("%s: %v", name, resp.Status)
		color.Yellow("%s: %v", name, string(resp_body))
		return fmt.Errorf(strings.TrimRight(string(resp_body), "\n"))
	}

	return nil
}

func checkFile(file string) (pass bool, err error) {
	f, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("file: %s not exist!", file)
		} else {
			return false, err
		}
	}

	if f.IsDir() {
		return false, fmt.Errorf("file: %s is a dir", file)
	}

	return true, nil
}
