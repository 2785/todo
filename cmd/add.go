package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/2785/todo/repositories"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a todo item",
	Long:  "TODO",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := repositories.GetR()
		if err != nil {
			return err
		}

		description := ""
		weight, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			weight = 10
			fmt.Fprintln(os.Stdout, "no weight provided, assuming weight of 10")
			description = strings.Join(args, " ")
		} else {
			description = strings.Join(args[1:], " ")
		}

		newTodo := repositories.Todo{
			Description: description,
			Weight:      int(weight),
		}

		added, err := r.AddTodo(cmd.Context(), newTodo)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, "added todo with id", added.ID.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
