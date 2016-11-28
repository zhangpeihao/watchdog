// Copyright © 2016 Zhang Peihao <zhangpeihao@gmail.com>
//

package cmd

import (
	"fmt"

	"github.com/zhangpeihao/watchdog/pkg/client"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status",
	Long: `
Status

Show all watching nginx statue pages.`,
	Run: func(cmd *cobra.Command, args []string) {
		status := client.GetStatus()
		for i, line := range status {
			fmt.Printf("%d - %s\n", i, line.URL)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
