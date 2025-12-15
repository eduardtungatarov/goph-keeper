package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/eduardtungatarov/goph-keeper/internal/client/handler"

	"github.com/spf13/cobra"
)

type Command struct {
	h *handler.Handler
}

func New(h *handler.Handler) *Command {
	return &Command{
		h: h,
	}
}

func (c *Command) LoginCmd() *cobra.Command {
	var login, password string

	cmd := &cobra.Command{
		Use:   "login [LOGIN]",
		Short: "Login to server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			login = args[0]
			c.h.HandleLogin(login, password)
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")

	return cmd
}

func (c *Command) RegisterCmd() *cobra.Command {
	var login, password string

	cmd := &cobra.Command{
		Use:   "register [LOGIN]",
		Short: "Register new user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			login = args[0]
			c.h.HandleRegister(login, password)
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")

	return cmd
}

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
			c.h.HandleCreate("pwd", data, title)
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
			c.h.HandleCreate("card", data, title)
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
			c.h.HandleCreate("binary", data, title)
		},
	}

	cmd.Flags().StringVarP(&data, "data", "d", "", "Data content")
	cmd.Flags().StringVarP(&title, "title", "T", "", "Item title")
	cmd.MarkFlagRequired("data")
	cmd.MarkFlagRequired("title")

	return cmd
}

func (c *Command) ReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [ID]",
		Short: "Read data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var id int64
			fmt.Sscanf(args[0], "%d", &id)
			c.h.HandleRead(id)
		},
	}
	return cmd
}

func (c *Command) DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [ID]",
		Short: "Delete data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var id int64
			fmt.Sscanf(args[0], "%d", &id)
			c.h.HandleDelete(id)
		},
	}
	return cmd
}

func (c *Command) ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all data items",
		Run: func(cmd *cobra.Command, args []string) {
			c.h.HandleList()
		},
	}

	return cmd
}
