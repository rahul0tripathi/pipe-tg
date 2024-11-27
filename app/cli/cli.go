package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rahul0tripathi/pipetg/internal/integrations/tg"
)

func RunAuth(client *tg.Client) error {
	reader := bufio.NewReader(os.Stdin)

	ctx := context.Background()

	err := client.WithUncheckedContext(ctx, func(_ctx context.Context) error {
		conn, err := client.GetTgConnFromCtx(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get telegram client: %w", err)
		}

		sessionData, err := client.GetSessionConfig()
		if err == nil && sessionData != "" {
			fmt.Println("session present", sessionData)
			return nil
		}

		fmt.Println("Sending authentication code...")
		if err := client.SendCode(_ctx, conn); err != nil {
			return fmt.Errorf("failed to send code: %w", err)
		}

		fmt.Print("Enter the authentication code: ")
		code, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read code: %w", err)
		}
		code = strings.TrimSpace(code)

		if err := client.AuthenticateWithCode(_ctx, code, conn); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		fmt.Println("Successfully authenticated!")

		sessionData, err = client.GetSessionConfig()
		if err != nil {
			return fmt.Errorf("failed to get session config: %w", err)
		}
		fmt.Printf("Session Config: %s\n", sessionData)

		return nil
	})

	return err
}
