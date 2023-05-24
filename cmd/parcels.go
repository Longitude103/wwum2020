package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var parcelDeleteCmd = &cobra.Command{
	Use:   "parcel-delete",
	Short: "Delete a parcel",
	Long:  `Delete a parcel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete parcel called")
		fmt.Printf("user inside parcel delete: %+v\n", user)
	},
}

func init() {
	rootCmd.AddCommand(parcelDeleteCmd)
}
