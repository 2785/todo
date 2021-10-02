package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/2785/todo/repositories"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type todos []repositories.Todo

func (t todos) String(i int) string {
	return t[i].Description
}

func (t todos) Len() int {
	return len(t)
}

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark tasks as finished",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := repositories.GetR()
		if err != nil {
			return err
		}
		var tds todos

		tds, err = repo.ListTodo(cmd.Context(), "weight", true, false)
		if err != nil {
			return err
		}

		if len(tds) == 0 {
			fmt.Fprintln(os.Stdout, "There are no pending todos")
		}

		if len(args) != 0 {
			searchPhrase := strings.Join(args, " ")
			res := fuzzy.FindFrom(searchPhrase, tds)

			if len(res) == 0 {
				fmt.Fprintf(os.Stdout, "Could not find match with search phrase %q, printing all\n", searchPhrase)
			}

			resTodos := make(todos, len(res))
			for i, v := range res {
				resTodos[i] = tds[v.Index]
			}

			tds = resTodos
		}

		promptMap := make(map[string]*repositories.Todo, len(tds))
		for i := range tds {
			if _, ok := promptMap[tds[i].Description]; ok {
				promptMap[tds[i].Description+"-dup"] = &tds[i]
			} else {
				promptMap[tds[i].Description] = &tds[i]
			}
		}

		options := make([]string, 0, len(promptMap))
		for k := range promptMap {
			options = append(options, k)
		}

		todosToMarkDone := []string{}

		err = survey.AskOne(&survey.MultiSelect{
			Message: "Which TODOs to close?",
			Options: options,
		}, &todosToMarkDone)

		if err != nil {
			if err == terminal.InterruptErr {
				fmt.Fprintln(os.Stdout, "prompt interrupped")
				return nil
			}
			return err
		}

		if len(todosToMarkDone) < 1 {
			fmt.Fprintln(os.Stdout, "nothing marked")
		}

		eg := &errgroup.Group{}
		for _, v := range todosToMarkDone {
			item := promptMap[v]
			eg.Go(func() error {
				_, err := repo.UpdateTodoStatus(cmd.Context(), *item, true)
				return err
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Marked %v items as done", len(todosToMarkDone))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
