package command

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eduardtungatarov/goph-keeper/internal/client/handler"

	"github.com/spf13/cobra"
)

// Command команды.
type Command struct {
	h *handler.Handler
}

// New конструктор команд.
func New(h *handler.Handler) *Command {
	return &Command{
		h: h,
	}
}

// LoginCmd команда входа.
func (c *Command) LoginCmd() *cobra.Command {
	var login, password string

	cmd := &cobra.Command{
		Use:   "login [LOGIN]",
		Short: "Login to server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			login = args[0]
			c.h.HandleLogin(cmd.Context(), login, password)
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")

	return cmd
}

// RegisterCmd команда регистрации.
func (c *Command) RegisterCmd() *cobra.Command {
	var login, password string

	cmd := &cobra.Command{
		Use:   "register [LOGIN]",
		Short: "Register new user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			login = args[0]
			c.h.HandleRegister(cmd.Context(), login, password)
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")

	return cmd
}

// CreateCmd команда сохранения данных на сервере.
func (c *Command) CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new data item",
	}

	cmd.AddCommand(
		c.createPwdCmd(),
		c.createCardCmd(),
		c.createBinaryCmd(),
	)

	return cmd
}

func (c *Command) createPwdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pwd",
		Short: "Create password item",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Title: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Login: ")
			login, _ := reader.ReadString('\n')
			login = strings.TrimSpace(login)

			fmt.Print("Password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			data := login + "|" + password
			c.h.HandleCreate(cmd.Context(), "pwd", data, title)
		},
	}

	return cmd
}

func (c *Command) createCardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Create card item",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Title: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Card number: ")
			number, _ := reader.ReadString('\n')
			number = strings.TrimSpace(number)

			fmt.Print("Expiry (MM/YY): ")
			exp, _ := reader.ReadString('\n')
			exp = strings.TrimSpace(exp)

			fmt.Print("CVV: ")
			cvv, _ := reader.ReadString('\n')
			cvv = strings.TrimSpace(cvv)

			data := number + "|" + exp + "|" + cvv
			c.h.HandleCreate(cmd.Context(), "card", data, title)
		},
	}

	return cmd
}

func (c *Command) createBinaryCmd() *cobra.Command {
	var data, title string

	cmd := &cobra.Command{
		Use:   "binary",
		Short: "Create binary item",
		Run: func(cmd *cobra.Command, args []string) {
			c.h.HandleCreate(cmd.Context(), "binary", data, title)
		},
	}

	cmd.Flags().StringVarP(&data, "data", "d", "", "Data content")
	cmd.Flags().StringVarP(&title, "title", "T", "", "Item title")
	cmd.MarkFlagRequired("data")
	cmd.MarkFlagRequired("title")

	return cmd
}

// ReadCmd команда чтения данных.
func (c *Command) ReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [ID]",
		Short: "Read data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			idStr := args[0]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Printf("Invalid ID '%s': %v\n", idStr, err)
				return
			}
			c.h.HandleRead(cmd.Context(), int64(id))
		},
	}
	return cmd
}

// DeleteCmd команда удаления данных.
func (c *Command) DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [ID]",
		Short: "Delete data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			idStr := args[0]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Printf("Invalid ID '%s': %v\n", idStr, err)
				return
			}
			c.h.HandleDelete(cmd.Context(), int64(id))
		},
	}
	return cmd
}

// ListCmd список данных пользователя.
func (c *Command) ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all data items",
		Run: func(cmd *cobra.Command, args []string) {
			c.h.HandleList(cmd.Context())
		},
	}

	return cmd
}
