package cmd

import (
	"fmt"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/save"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies planned changes on Azure cloud",
	Long: `Applies planned changes on Azure cloud. 

Using plan file created with 'plan' command this command performs actual 'terraform apply' operation. This command
performs following steps: 
 - validates presence of config and module state files
 - checks that module status is either 'Initialized' or 'Destroyed'
 - performs 'terraform apply' operation using existing plan file
 - updates module state file with applied config
 - saves terraform output to module state file. 

This command should always be preceded by 'plan' command.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

		accessKeyId = viper.GetString("access_key_id")
		secretAccessKey = viper.GetString("secret_access_key")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("apply called")

		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		config, state, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("checkAndLoad failed")
		}

		err = showModulePlan(config, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("showModulePlan failed")
		}
		output, err := terraformApply()
		if err != nil {
			logger.Fatal().Err(err).Msgf("registered following output: \n%s\n", output)
		}

		state.AwsBI.Config = config
		state.AwsBI.Status = st.Applied

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}

		terraformOutputMap, err := getTerraformOutputMap()
		if err != nil {
			logger.Fatal().Err(err).Msg("getTerraformOutputMap failed")
		}
		state.AwsBI.Output = produceOutput(terraformOutputMap)

		logger.Debug().Msg("save state")
		err = save.State(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}

		bytes, err := state.Marshal()
		if err != nil {
			logger.Fatal().Err(err).Msg("state.Marshal failed")
		}
		logger.Info().Msg(string(bytes))
		fmt.Println("State after apply: \n" + string(bytes))

		msg, err := count(output)
		if err != nil {
			logger.Fatal().Err(err).Msg("count failed")
		}
		logger.Info().Msg("Performed following changes: " + msg)
		fmt.Println("Performed following changes: \n\t" + msg)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().String("ACCESS_KEY_ID", "", "AWS access key ID")
	applyCmd.Flags().String("SECRET_ACCESS_KEY", "", "AWS secret access key")
}
