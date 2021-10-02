package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/2785/todo/repositories"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo app for the clinically insane",
	Args:  cobra.ArbitraryArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch viper.GetString("storage_mode") {
		case "":
			return errors.New("'storage_mode' must be set")
		case "file":
			if viper.GetString("file_storage_dir") == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}

				defaultStoragePath := path.Join(home, ".todo", "flatfile-storage")
				if _, err := os.Stat(defaultStoragePath); os.IsNotExist(err) {
					err := os.MkdirAll(defaultStoragePath, os.ModePerm)
					if err != nil {
						return err
					}
				}

				viper.Set("file_storage_dir", defaultStoragePath)
			}

			flatFileRepo, err := repositories.NewFlatFile(viper.GetString("file_storage_dir"))
			if err != nil {
				return err
			}
			repositories.SetR(flatFileRepo)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			return addCmd.RunE(cmd, args)
		}
		return cmd.Help()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")
	viper.SetDefault("storage_mode", "file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".todo")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
