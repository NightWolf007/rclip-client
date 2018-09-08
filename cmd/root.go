package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NightWolf007/rclip-client/listeners"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rclip-client",
	Short: "Client for RClip",
	Long:  `Client for RClip`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if viper.GetBool("debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting listeners...")

		servers := viper.GetStringSlice("servers")
		timeout := viper.GetDuration("timeout")
		updateDelay := viper.GetDuration("update_delay")
		recoverDelay := viper.GetDuration("recover_delay")

		var wg sync.WaitGroup
		wg.Add(2 * len(servers))

		for _, addr := range servers {
			go func(addr string) {
				defer wg.Done()

				for {
					err := listeners.RunRemoteListener(addr, timeout)
					if err != nil {
						log.Error().Err(err).Msg("Remote listener exited")
					}
					time.Sleep(recoverDelay)
				}
			}(addr)
			go func(addr string) {
				defer wg.Done()

				for {
					err := listeners.RunLocalListener(addr, timeout, updateDelay)
					if err != nil {
						log.Error().Err(err).Msg("Local listener exited")
					}
					time.Sleep(recoverDelay)
				}
			}(addr)
		}

		wg.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "",
		"config file (default is $HOME/.rclipd-client.yaml)",
	)

	rootCmd.PersistentFlags().BoolP(
		"debug", "d", false,
		"Enable debug output",
	)
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringSliceP(
		"servers", "s", []string{"localhost:9889"},
		"RClip servers",
	)
	viper.BindPFlag("servers", rootCmd.PersistentFlags().Lookup("servers"))

	rootCmd.PersistentFlags().DurationP(
		"timeout", "t", 5*time.Second,
		"RClip connection timeout",
	)
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	rootCmd.PersistentFlags().DurationP(
		"update-delay", "u", 2*time.Second,
		"Delay between checks of local clipboard",
	)
	viper.BindPFlag("update_delay", rootCmd.PersistentFlags().Lookup("update-delay"))

	rootCmd.PersistentFlags().DurationP(
		"recover-delay", "r", 5*time.Second,
		"Delay before recover connection after failure",
	)
	viper.BindPFlag("recover_delay", rootCmd.PersistentFlags().Lookup("recover-delay"))
}

// initConfig reads in config file
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rclipd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rclipd-client")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
