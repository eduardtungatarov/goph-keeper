package command

import (
	"fmt"

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
	var dataType, data, title string

	cmd := &cobra.Command{
		Use:   "create [TYPE]",
		Short: "Create new data item",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataType = args[0]
			c.h.HandleCreate(dataType, data, title)
		},
	}

	cmd.Flags().StringVarP(&dataType, "type", "t", "", "Data type: pwd, card, binary")
	cmd.Flags().StringVarP(&data, "data", "d", "", "Data content")
	cmd.Flags().StringVarP(&title, "title", "T", "", "Item title")
	cmd.MarkFlagRequired("data")
	cmd.MarkFlagRequired("title")

	return cmd
}

func (h *Command) ReadCmd() *cobra.Command {
	var id int64

	cmd := &cobra.Command{
		Use:   "read [ID]",
		Short: "Read data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Sscanf(args[0], "%d", &id)
			// h.Commandead(id)
		},
	}

	cmd.Flags().Int64VarP(&id, "id", "i", 0, "Data item ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func (h *Command) DeleteCmd() *cobra.Command {
	var id int64

	cmd := &cobra.Command{
		Use:   "delete [ID]",
		Short: "Delete data item by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Sscanf(args[0], "%d", &id)
			// h.handleDelete(id)
		},
	}

	cmd.Flags().Int64VarP(&id, "id", "i", 0, "Data item ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func (h *Command) ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all data items",
		Run: func(cmd *cobra.Command, args []string) {
			// h.handleList()
		},
	}

	return cmd
}
