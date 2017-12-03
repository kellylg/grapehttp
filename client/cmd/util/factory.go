package util

import (
	"encoding/base64"
	"flag"
	"github.com/parnurzeal/gorequest"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"grapehttp/pkg/homedir"
)

const (
	RecommendedHomeDir = ".grape"
)

type Factory struct {
	flags    *pflag.FlagSet
	Server   string
	Timeout  int
	Username string
	Password string
	Cool     bool
}

func NewFactory() Factory {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	f := Factory{
		flags:    flags,
		Server:   viper.GetString("server"),
		Timeout:  viper.GetInt("timeout"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Cool:     viper.GetBool("cool"),
	}

	return f
}

func (f *Factory) FlagSet() *pflag.FlagSet {
	return f.flags
}

// TODO: We need to filter out stuff like secrets.
func (f *Factory) Command() string {
	if len(os.Args) == 0 {
		return ""
	}
	base := filepath.Base(os.Args[0])
	args := append([]string{base}, os.Args[1:]...)
	return strings.Join(args, " ")
}

func (f *Factory) BindFlags(flags *pflag.FlagSet) {
	// Merge factory's flags
	flags.AddFlagSet(f.flags)
}

func (f *Factory) BindExternalFlags(flags *pflag.FlagSet) {
	// any flags defined by external projects (not part of pflags)
	flags.AddGoFlagSet(flag.CommandLine)
}

func (f *Factory) Auth() string {
	return base64.StdEncoding.EncodeToString([]byte(f.Username + ":" + f.Password))
}

func (f *Factory) Gorequest() *gorequest.SuperAgent {
	request := gorequest.New()
	request.DoNotClearSuperAgent = true
	return request.Timeout(time.Duration(f.Timeout)*time.Second).
		Set("Content-Type", "application/json").
		Set("Username", "cc")
}

func init() {
	viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config") // no need to include file extension

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
