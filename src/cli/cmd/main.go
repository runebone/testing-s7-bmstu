package cmd

import (
	"cli/internal/adapter/logger"
	"cli/internal/adapter/service/http"
	"cli/internal/config"
	"cli/internal/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	_ "time/tzdata"

	v1 "cli/internal/usecase/v1"

	"github.com/spf13/cobra"
)

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Couldn't set timezone: %v", err)
	}
	time.Local = loc
}

func readTokens(tokens *dto.Tokens, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", filePath)
		return err
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(tokens); err != nil {
		fmt.Printf("Error parsing file: %s\n", filePath)
		return err
	}

	return nil
}

func saveTokens(tokens *dto.Tokens, filePath string) error {
	tokensJson, err := json.Marshal(tokens)
	if err != nil {
		fmt.Println("Error marshalling tokens")
		return err
	}

	err = os.WriteFile(filePath, tokensJson, 0644)
	if err != nil {
		fmt.Printf("Error saving tokens to file: %s\n", filePath)
	}

	return nil
}

func main() {
	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	var Tokens dto.Tokens
	err = readTokens(&Tokens, config.Client.TokensPath)
	if err != nil {
		return
	}

	logger := logger.NewZapLogger(config.Aggregator.Log)
	ac := config.Aggregator
	svc := http.NewAggregatorService(ac.BaseURL, 5*time.Second, logger)

	client := v1.NewClientUseCase(svc)

	rootCmd := &cobra.Command{Use: "todo"}

	// Register command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "register [username] [email] [password]",
		Short: "Register a new user",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := client.Register(context.Background(), args[0], args[1], args[2])
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			saveTokens(tokens, config.Client.TokensPath)
		},
	})

	// Login command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "login [email] [password]",
		Short: "Login a user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := client.Login(context.Background(), args[0], args[1])
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			saveTokens(tokens, config.Client.TokensPath)
		},
	})

	// Logout command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Logout the user",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Logout(context.Background(), Tokens.RefreshToken)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Println("logged out")
		},
	})

	// Create board
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create board [title]",
		Short: "Create a new board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.CreateBoard(ctx, args[0])
		},
	})

	// Create column
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create column [board_id] [title]",
		Short: "Create a new column in a board",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.CreateColumn(ctx, args[0], args[1])
		},
	})

	// Create card
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create card [column_id] [title] [description]",
		Short: "Create a new card in a column",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			var description string
			if len(args) == 3 {
				description = args[2]
			} else {
				description = ""
			}
			client.CreateCard(ctx, args[0], args[1], description)
		},
	})

	// Show boards
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show boards",
		Short: "Show all boards",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.ShowBoards(ctx)
		},
	})

	// Show board
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show board [board_id]",
		Short: "Show a board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.ShowBoard(ctx, args[0])
		},
	})

	// Show column
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show column [column_id]",
		Short: "Show a column",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.ShowColumn(ctx, args[0])
		},
	})

	// Show card
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show card [card_id]",
		Short: "Show a card",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.ShowCard(ctx, args[0])
		},
	})

	// Update board title
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update board [board_id] title [new_title]",
		Short: "Update board title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.UpdateBoard(ctx, args[0], args[1])
		},
	})

	// Update column title
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update column [column_id] title [new_title]",
		Short: "Update column title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.UpdateColumn(ctx, args[0], args[1])
		},
	})

	// Update card title
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update card [card_id] title [new_title]",
		Short: "Update card title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.UpdateCardTitle(ctx, args[0], args[1])
		},
	})

	// Update card description
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update card [card_id] description [new_description]",
		Short: "Update card description",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.UpdateCardDescription(ctx, args[0], args[1])
		},
	})

	// Delete board
	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete board [board_id]",
		Short: "Delete a board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.DeleteBoard(ctx, args[0])
		},
	})

	// Delete column
	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete column [column_id]",
		Short: "Delete a column",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.DeleteColumn(ctx, args[0])
		},
	})

	// Delete card
	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete card [card_id]",
		Short: "Delete a card",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.DeleteCard(ctx, args[0])
		},
	})

	// Stats
	rootCmd.AddCommand(&cobra.Command{
		Use:   "stats from [DD-MM-YYYY] to [DD-MM-YYYY]",
		Short: "Show stats for a time period",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(context.Background(), "tokens", Tokens)
			client.Stats(ctx, args[0], args[1])
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
