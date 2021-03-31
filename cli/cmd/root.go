package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/masnax/rest-api/cli/cmd/add"
	"github.com/masnax/rest-api/cli/cmd/delete"
	"github.com/masnax/rest-api/cli/cmd/edit"
	"github.com/masnax/rest-api/cli/cmd/list"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const URL string = "http://localhost:8080/"

//flags
var (
	filterFlag      string
	collectionFlag  string
	bookFlag        string
	titleFlag       string
	authorFlag      string
	dateFlag        string
	editionFlag     int
	descriptionFlag string
	genreFlag       string
)

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
		filter, ok := parseFilter(cmd, filterFlag)
		if ok {
			argPath += filter
			header, data := list.GetBookList(URL, "books", argPath)
			renderTable(header, data)
		}
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

var cmdEditBook = &cobra.Command{
	Use:   "edit id",
	Short: "Update book with id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		edit.EditBook(URL, "books", args[0], titleFlag, authorFlag, dateFlag,
			editionFlag, descriptionFlag, genreFlag)
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
			filter, ok := parseFilter(cmd, filterFlag)
			if ok {
				argPath += filter
				header, data = list.GetBookList(URL, "collections", argPath)
			}
		} else if len(bookFlag) > 0 {
			argPath += "/book/" + bookFlag
			header, data = list.GetCollectionList(URL, "collections", argPath)
		} else {
			header, data = list.GetCollectionStatList(URL, "collections", argPath)
		}
		renderTable(header, data)
	},
}

var cmdAddCollection = &cobra.Command{
	Use:   "new name",
	Short: "Add a new collection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		add.AddNewCollection(URL, "collections/manage", args)
	},
}

var cmdEditCollection = &cobra.Command{
	Use:   "edit id name",
	Short: "Update collection name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		edit.EditCollection(URL, "collections/manage", args[0], args[1])
	},
}

var cmdDelCollection = &cobra.Command{
	Use:     "delete id",
	Aliases: []string{"rm"},
	Short:   "Delete collection with id",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		delete.DelBook(URL, "collections/manage", args)
	},
}

var cmdAddToCollection = &cobra.Command{
	Use:   "add book_id collection_id",
	Short: "Add a book to a collection",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		add.AddToCollection(URL, "collections", args[0], args[1])
	},
}

var cmdRemoveFromCollection = &cobra.Command{
	Use:   "drop book_id collection_id",
	Short: "Drop a book from a collection",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		delete.RemoveFromCollection(URL, "collections", args[0], args[1])
	},
}

func parseArgs(args []string) string {
	argPath := ""
	for _, a := range args {
		argPath += "/" + a
	}
	return argPath
}

func parseFilter(cmd *cobra.Command, filter string) (string, bool) {
	if len(filter) == 0 {
		return "", true
	}
	parts := strings.Split(filter, " ")
	out := "?filter="
	if len(parts) < 3 {
		log.Println(cmd.Flag("filter").Usage)
		return "", false
	}
	out += strings.ReplaceAll(filter, " ", "+")
	return out, true
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
	cmdEditBook.Flags().StringVar(&titleFlag, "title", "", "book title")
	cmdEditBook.Flags().StringVar(&authorFlag, "author", "", "book author")
	cmdEditBook.Flags().StringVar(&dateFlag, "published", "", "book publish date")
	cmdEditBook.Flags().IntVar(&editionFlag, "edition", 0, "book edition")
	cmdEditBook.Flags().StringVar(&descriptionFlag, "description", "", "book description")
	cmdEditBook.Flags().StringVar(&genreFlag, "genre", "", "book genre")

	cmdListCollections.Flags().StringVarP(&filterFlag, "filter", "f", "",
		"'--filter' format: \"key [eq,ne,lt,gt,le,ge] value\"")
	cmdListBooks.Flags().StringVarP(&filterFlag, "filter", "f", "",
		"'--filter' format: \"key [eq,ne,lt,gt,le,ge] value\"")

	rootCmd.AddCommand(cmdCollections)
	cmdCollections.AddCommand(cmdListCollections)
	cmdCollections.AddCommand(cmdAddCollection)
	cmdCollections.AddCommand(cmdEditCollection)
	cmdCollections.AddCommand(cmdDelCollection)
	cmdCollections.AddCommand(cmdRemoveFromCollection)
	cmdCollections.AddCommand(cmdAddToCollection)

	rootCmd.AddCommand(cmdListBooks)
	rootCmd.AddCommand(cmdAddBook)
	rootCmd.AddCommand(cmdDelBook)
	rootCmd.AddCommand(cmdEditBook)

	return rootCmd.Execute()
}
