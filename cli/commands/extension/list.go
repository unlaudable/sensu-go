package extension

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/sensu/sensu-go/cli"
	"github.com/sensu/sensu-go/cli/client"
	"github.com/sensu/sensu-go/cli/commands/helpers"
	"github.com/sensu/sensu-go/cli/elements/table"
	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

func ListCommand(cli *cli.SensuCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list extensions",
		RunE:  runList(cli.Config.Format(), cli.Client, cli.Config.Organization(), cli.Config.Format()),
	}
	helpers.AddFormatFlag(cmd.Flags())
	return cmd
}

func runList(config string, client client.APIClient, org, format string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			_ = cmd.Help()
			return errors.New("invalid arguments received")
		}
		extensions, err := client.ListExtensions(org)
		if err != nil {
			return err
		}

		// Print the results based on the user preferences
		resources := []types.Resource{}
		for i := range extensions {
			resources = append(resources, &extensions[i])
		}
		return helpers.Print(cmd, config, printToTable, resources, extensions)
	}
}

func printToTable(results interface{}, writer io.Writer) {
	table := table.New([]*table.Column{
		{
			Title:       "Name",
			ColumnStyle: table.PrimaryTextStyle,
			CellTransformer: func(data interface{}) string {
				asset, _ := data.(types.Asset)
				return asset.Name
			},
		},
		{
			Title: "URL",
			CellTransformer: func(data interface{}) string {
				asset, _ := data.(types.Asset)
				u, err := url.Parse(asset.URL)
				if err != nil {
					return ""
				}

				_, file := path.Split(u.EscapedPath())
				return fmt.Sprintf(
					"//%s/.../%s",
					u.Hostname(),
					file,
				)
			},
		},
	})

	table.Render(writer, results)
}
