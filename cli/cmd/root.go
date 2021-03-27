package cmd

import (
	"os"

	"github.com/masnax/rest-api/cli/cmd/add"
	"github.com/masnax/rest-api/cli/cmd/delete"
	"github.com/masnax/rest-api/cli/cmd/edit"
	"github.com/masnax/rest-api/cli/cmd/list"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const URL string = "http://localhost:8080/"

var collectionFlag string
var bookFlag string
var rootCmd = &cobra.Command{
	Use:   "bmc",
	Short: "bmc - book manager",
	Long:  "bmc is a book manager for managing books and collections of books",
}

var cmdListBooks = &cobra.Command{
	Use:     "list [id] [flags]",
	Aliases: []string{"ls"},
	Short:   "List books",
	Run: func(cmd *cobra.Command, args []string) {
		argPath := parseArgs(args)
		header, data := list.GetBookList(URL, "books", argPath)
		renderTable(header, data)
	},
}

var cmdAddBook = &cobra.Command{
	Use:   "add title author published edition description genre",
	Short: "Add a new book with the given attributes",
	Args:  cobra.ExactArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		add.AddBook(URL, "books", args)
	},
}

var cmdDelBook = &cobra.Command{
	Use:     "delete id",
	Aliases: []string{"rm"},
	Short:   "Delete book with id",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		delete.DelBook(URL, "books", args)
	},
}

var cmdEditbook = &cobra.Command{
	Use:   "edit id",
	Short: "Update book with id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		edit.EditBook(URL, "books", args)
	},
}

var cmdCollections = &cobra.Command{
	Use:     "collection [command]",
	Aliases: []string{"col"},
	Short:   "Manage collections of books",
	Long: `Manage collections of books:
	add, update, show, and delete collections and their associated books`,
}

var cmdListCollections = &cobra.Command{
	Use:     "list [command]",
	Aliases: []string{"ls"},
	Short:   "List collections and their books",
	Run: func(cmd *cobra.Command, args []string) {
		argPath := parseArgs(args)
		var header []string
		var data [][]string
		if len(collectionFlag) > 0 {
			argPath += "/collection/" + collectionFlag
			header, data = list.GetBookList(URL, "collections", argPath)
		} else if len(bookFlag) > 0 {
			argPath += "/book/" + bookFlag
			header, data = list.GetCollectionList(URL, "collections", argPath)
		} else {
			header, data = list.GetCollectionStatList(URL, "collections", argPath)
		}
		renderTable(header, data)
	},
}

func parseArgs(args []string) string {
	argPath := ""
	for _, a := range args {
		argPath += "/" + a
	}
	return argPath
}

func renderTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}

func Execute() error {
	//var rootCmd = &cobra.Command{Use: "bmc"}
	cmdListCollections.Flags().StringVar(&collectionFlag, "name", "",
		"shows all books for a given collection name")
	cmdListCollections.Flags().StringVar(&bookFlag, "bid", "",
		"shows all collections for a given book id")
	rootCmd.AddCommand(cmdListBooks)
	rootCmd.AddCommand(cmdCollections)
	rootCmd.AddCommand(cmdAddBook)
	rootCmd.AddCommand(cmdDelBook)
	cmdCollections.AddCommand(cmdListCollections)

	return rootCmd.Execute()
}
