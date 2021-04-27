package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	doDestroy bool
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "performs module plan operation",
	Long: `Performs module plan operation. 

There is two steps performed currently: 
 - simulate how module state will change after apply
 - perform 'terraform plan' operation (simulate what resources will be installed). 

Predicted module state is not being recorded but terraform plan file is being created. It means that in consecutive 
invoking of 'apply' command, created plan file would be used.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

		accessKeyId = viper.GetString("access_key_id")
		secretAccessKey = viper.GetString("secret_access_key")
		doDestroy = viper.GetBool("destroy")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("plan called")
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		config, state, err := checkAndLoad(stateFilePath, configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("checkAndLoad failed")
		}

		err = templateTfVars(config)
		if err != nil {
			logger.Fatal().Err(err).Msg("templateTfVars failed")
		}
		if !doDestroy {
			err = showModulePlan(config, state)
			if err != nil {
				logger.Fatal().Err(err).Msg("showModulePlan failed")
			}
			output, err := terraformPlan()
			if err != nil {
				logger.Fatal().Err(err).Msg("terraform plan failed")
			}
			msg, err := count(output)
			if err != nil {
				logger.Fatal().Err(err).Msg("count failed")
			}
			logger.Info().Msg("Will perform following changes: " + msg)
			fmt.Println("Will perform following changes: \n\t" + msg)
		} else {
			output, err := terraformPlanDestroy()
			if err != nil {
				logger.Fatal().Err(err).Msg("terraform plan destroy failed")
			}
			msg, err := count(output)
			if err != nil {
				logger.Fatal().Err(err).Msg("count failed")
			}
			logger.Info().Msg("Will perform following changes: " + msg)
			fmt.Println("Will perform following changes: \n\t" + msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.Flags().String("ACCESS_KEY_ID", "", "AWS access key ID")
	planCmd.Flags().String("SECRET_ACCESS_KEY", "", "AWS secret access key")
	planCmd.Flags().Bool("destroy", false, "prepare plan for destroy")
}
