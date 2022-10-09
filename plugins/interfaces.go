package plugins

import "github.com/spf13/cobra"

type Plugin interface {
	Init(cmd *cobra.Command) error
	Update() error
	GetName() string
}
