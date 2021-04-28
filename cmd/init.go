package cmd

import (
	"errors"
	"fmt"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/load"
	"github.com/epiphany-platform/e-structures/utils/save"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
)

var (
	name       string
	vmsRsaPath string
	region     string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes module configuration file",
	Long:  `Initializes module configuration file (in ` + filepath.Join(defaultSharedDirectory, moduleShortName, configFileName) + `/). `,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("initialization error occurred")
		}

		name = viper.GetString("name")
		vmsRsaPath = viper.GetString("vms_rsa")
		region = viper.GetString("region")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("init called")
		moduleDirectoryPath := filepath.Join(SharedDirectory, moduleShortName)
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		logger.Debug().Msg("ensure directories")
		err := ensureDirectory(moduleDirectoryPath)
		if err != nil {
			logger.Fatal().Err(err).Msg("ensureDirectory failed")
		}
		logger.Debug().Msg("load state file")
		state, err := load.State(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("loadState failed")
		}
		logger.Debug().Msg("load config file")
		config, err := load.AwsBIConfig(configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("loadConfig failed")
		}

		if state.AwsBI != nil {
			if !reflect.DeepEqual(state.AwsBI, &st.AwsBIState{}) && state.AwsBI.Status != st.Initialized && state.AwsBI.Status != st.Destroyed {
				logger.Fatal().Err(errors.New(string("unexpected state: " + state.AwsBI.Status))).Msg("incorrect state")
			}
		} else {
			state.AwsBI = &st.AwsBIState{}
		}

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}
		logger.Debug().Msg("backup config file")
		err = backupFile(configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}

		config.Params.Name = to.StrPtr(name)
		config.Params.Region = to.StrPtr(region)
		config.Params.RsaPublicKeyPath = to.StrPtr(filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)))

		state.AwsBI.Status = st.Initialized

		logger.Debug().Msg("save config")
		err = save.AwsBIConfig(configFilePath, config)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveConfig failed")
		}
		logger.Debug().Msg("save state")
		err = save.State(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}

		bytes, err := config.Marshal()
		if err != nil {
			logger.Fatal().Err(err).Msg("config.Marshal failed")
		}
		logger.Info().Msg(string(bytes))
		fmt.Println("Initialized config: \n" + string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().String("name", "epiphany", "prefix given to all resources created") //TODO rename to prefix
	initCmd.Flags().String("region", "eu-central-1", "region used to initialize config file")
	initCmd.Flags().String("vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines")
}
