package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the User Information",
	Long:  `Initializes the User Information so that you can log in. Must be run before any other commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		newUser := User{
			Username: user,
			Password: password,
		}

		// write to file and save in dir
		newUserBytes, err := json.Marshal(newUser)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(".configjson", newUserBytes, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Successfully initialized")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	// Here you will define your flags and configuration settings.
	initCmd.Flags().StringP("username", "u", "", "Username for login to database")
	initCmd.Flags().StringP("password", "p", "", "Password for login to database")
}
