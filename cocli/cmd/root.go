// Copyright 2021-2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/veraison/apiclient/auth"
)

var (
	cfgFile string
	fs      = afero.NewOsFs()

	cliConfig  = &ClientConfig{}
	authMethod = auth.MethodPassthrough
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "cocli",
	Short:         "CoRIM & CoMID swiss-army knife",
	Version:       "0.0.1",
	SilenceUsage:  true,
	SilenceErrors: true,
}

type ClientConfig struct {
	Auth auth.IAuthenticator
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/cocli/config.yaml)")
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	v, err := readConfig(cfgFile)
	cobra.CheckErr(err)

	err = authMethod.Set(v.GetString("auth"))
	cobra.CheckErr(err)

	switch authMethod {
	case auth.MethodPassthrough:
		cliConfig.Auth = &auth.NullAuthenticator{}
	case auth.MethodBasic:
		cliConfig.Auth = &auth.BasicAuthenticator{}
		err = cliConfig.Auth.Configure(map[string]interface{}{
			"username": v.GetString("username"),
			"password": v.GetString("password"),
		})
		cobra.CheckErr(err)
	case auth.MethodOauth2:
		cliConfig.Auth = &auth.Oauth2Authenticator{}
		err = cliConfig.Auth.Configure(map[string]interface{}{
			"client_id":     v.GetString("client_id"),
			"client_secret": v.GetString("client_secret"),
			"token_url":     v.GetString("token_url"),
			"username":      v.GetString("username"),
			"password":      v.GetString("password"),
			"ca_certs":      v.GetStringSlice("ca_cert"),
		})
		cobra.CheckErr(err)
	default:
		// Should never get here as authMethod value is set via
		// Method.Set(), which ensures that it's one of the above.
		panic(fmt.Sprintf("unknown auth method: %q", authMethod))
	}
}

func readConfig(path string) (*viper.Viper, error) {
	v := viper.GetViper()
	if path != "" {
		v.SetConfigFile(path)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		userConfigDir, err := os.UserConfigDir()
		if err == nil {
			v.AddConfigPath(filepath.Join(userConfigDir, "cocli"))
		}
		v.AddConfigPath(wd)
		v.SetConfigType("yaml")
		v.SetConfigName("config")
	}

	v.SetEnvPrefix("cocli")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		err = nil
	}

	return v, err
}
