package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "rcl",
		Short: "recall",
		RunE: func(cmd *cobra.Command, args []string) error {
			readFile, err := os.Open(os.Getenv("HISTFILE"))
			if err != nil {
				return err
			}
			defer readFile.Close()

			fileScanner := bufio.NewScanner(readFile)
			fileScanner.Split(bufio.ScanLines)

			var result []string
			for fileScanner.Scan() {
				result = append(result, fileScanner.Text())
			}

			for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
				result[i], result[j] = result[j], result[i]
			}

			fmt.Println(result)
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	errGrp, errctx := errgroup.WithContext(ctx)
	errGrp.Go(func() error {
		defer cancel()
		if err := rootCmd.ExecuteContext(errctx); err != nil {
			return err
		}
		return nil
	})
	return errGrp.Wait()
}
