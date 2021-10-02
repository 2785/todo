package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/2785/todo/repositories"
	"github.com/spf13/cobra"
)

var (
	listShowAll  bool
	listRankBy   string
	listRankDesc bool
)

var (
	todoTmpl = func() *template.Template {
		tmpl, err := template.New("todo").Parse(`W: {{.Weight}} - added {{ .CreatedAt.Format "Jan 02 15:04 Mon"}}
ID: {{.ID}}
Desc: {{.Description}}`)
		if err != nil {
			panic(err)
		}
		return tmpl
	}()
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "list the todos",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := repositories.GetR()
		if err != nil {
			return err
		}

		todos, err := repo.ListTodo(cmd.Context(), listRankBy, listRankDesc, listShowAll)
		if err != nil {
			return err
		}

		if len(todos) == 0 {
			fmt.Fprintln(os.Stdout, "no todos found")
		}

		for i, v := range todos {
			err := todoTmpl.Execute(os.Stdout, v)
			if err != nil {
				return err
			}
			if i != len(todos)-1 {
				fmt.Fprintln(os.Stdout, "")
				fmt.Fprintln(os.Stdout, strings.Repeat("-", 10))
			}
		}

		fmt.Fprint(os.Stdout, "\n\n")

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&listShowAll, "all", false, "whether to show completed todos")
	listCmd.Flags().BoolVar(&listRankDesc, "desc", false, "whether to sort desc")
	listCmd.Flags().StringVar(&listRankBy, "by", "weight", "what to sort by")

	rootCmd.AddCommand(listCmd)
}
