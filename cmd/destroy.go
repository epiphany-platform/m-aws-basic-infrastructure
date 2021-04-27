package cmd

import (
	"fmt"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/save"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "performs planed destroy operation",
	Long: `Performs planed destroy operation.

This command uses terraform plan file prepared with previously called 'plan --destroy' command. Before destroy it checks
if module status is 'Applied'. `,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("initialization error occurred")
		}

		accessKeyId = viper.GetString("access_key_id")
		secretAccessKey = viper.GetString("secret_access_key")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("destroy called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		config, state, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("checkAndLoad failed")
		}

		if state.AwsBI == nil {
			state.AwsBI = &st.AwsBIState{}
		}

		err = templateTfVars(config)
		if err != nil {
			logger.Fatal().Err(err).Msg("templateTfVars failed")
		}
		output, err := terraformDestroy()
		if err != nil {
			logger.Fatal().Err(err).Msgf("registered following output: \n%s\n", output)
		}
		msg, err := count(output)
		if err != nil {
			logger.Fatal().Err(err).Msg("count failed")
		}
		logger.Info().Msg("Performed following changes: " + msg)
		fmt.Println("Performed following changes: \n\t" + msg)
		state = updateStateAfterDestroy(state)
		err = save.State(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().String("ACCESS_KEY_ID", "", "AWS access key ID")
	destroyCmd.Flags().String("SECRET_ACCESS_KEY", "", "AWS secret access key")
}
