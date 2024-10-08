package cmd

import (
	"cli/internal/adapter/logger"
	"cli/internal/adapter/service/http"
	"cli/internal/config"
	"cli/internal/dto"
	"context"
	"fmt"
	"log"
	"time"

	v1 "cli/internal/usecase/v1"

	"github.com/spf13/cobra"
)

func main() {
	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
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
			client.Register(context.Background(), args[0], args[1], args[2])
		},
	})

	// Login command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "login [email] [password]",
		Short: "Login a user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client.Login(context.Background(), args[0], args[1])
		},
	})

	// Logout command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Logout the user",
		Run: func(cmd *cobra.Command, args []string) {
			client.Logout(context.Background(), "") // Assuming refresh token management
		},
	})

	// Create board
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create board [title]",
		Short: "Create a new board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			board := dto.Board{Title: args[0]}
			client.CreateBoard(context.Background(), board)
		},
	})

	// Create column
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create column [board_id] [title]",
		Short: "Create a new column in a board",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			column := dto.Column{BoardID: args[0], Title: args[1]}
			client.CreateColumn(context.Background(), column)
		},
	})

	// Create card
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create card [board_id] [column_id] [title] [description]",
		Short: "Create a new card in a column",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			card := dto.Card{
				BoardID:     args[0],
				ColumnID:    args[1],
				Title:       args[2],
				Description: "",
			}
			if len(args) == 4 {
				card.Description = args[3]
			}
			client.CreateCard(context.Background(), card)
		},
	})

	// Show boards
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show boards",
		Short: "Show all boards",
		Run: func(cmd *cobra.Command, args []string) {
			client.ShowBoards(context.Background())
		},
	})

	// Show board
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show board [board_id]",
		Short: "Show a board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client.ShowBoard(context.Background(), args[0])
		},
	})

	// Update board title
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update board [board_id] title [new_title]",
		Short: "Update board title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			board := dto.Board{ID: args[0], Title: args[1]}
			client.UpdateBoard(context.Background(), &board)
		},
	})

	// Delete board
	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete board [board_id]",
		Short: "Delete a board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client.DeleteBoard(context.Background(), args[0])
		},
	})

	// Stats
	rootCmd.AddCommand(&cobra.Command{
		Use:   "stats from [DD-MM-YYYY] to [DD-MM-YYYY]",
		Short: "Show stats for a time period",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client.Stats(context.Background(), args[0], args[1])
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Stubbed Client implementation for testing
func NewClient() Client {
	return &StubClient{}
}

// StubClient implements the Client interface with dummy methods
type StubClient struct{}

func (s *StubClient) Register(ctx context.Context, username, email, password string) {
	fmt.Println("Registered:", username, email)
}

func (s *StubClient) Login(ctx context.Context, email, password string) {
	fmt.Println("Logged in:", email)
}

func (s *StubClient) Logout(ctx context.Context, refreshToken string) {
	fmt.Println("Logged out")
}

// Implement other methods similarly...
