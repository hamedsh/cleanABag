package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Strubbl/wallabago/v7"
	"github.com/spf13/cobra"
)

var (
	cfgUnread  bool
	cfgStarred bool
	cfgDate    string
	cfgDelete  bool
	cfgForce   bool
)

func init() {
	pruneCmd.PersistentFlags().BoolVarP(&cfgUnread, "unread", "u", false, "Include unread entries for deletion. False will prevent unread articles from being deleted")
	pruneCmd.Flags().BoolVarP(&cfgStarred, "starred", "s", false, "Include starred entry in deletion. False will prevent starred article to be deleted.")
	pruneCmd.Flags().StringVarP(&cfgDate, "date", "d", "", "Articles older than this date will be removed if they match the archived/starred flags, format \"YYYY-MM-DDTHH-MM\".")
	pruneCmd.Flags().BoolVar(&cfgDelete, "delete", false, "Delete articles. Without this flag, it will only do a dry run.")
	pruneCmd.Flags().BoolVar(&cfgForce, "force", false, "Delete_article, even its archived.")

	pruneCmd.MarkFlagRequired("date")

	rootCmd.AddCommand(pruneCmd)
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete old article from wallabag",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running prune command")

		home, _ := os.UserHomeDir()
		configJSON := home + "/.config/cleanABag/credentials.json"
		if cfgFile != "" {
			configJSON = cfgFile
		}
		err := wallabago.ReadConfig(configJSON)
		if err != nil {
			fmt.Println("Error reading credentials config file")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		baseDate, err := time.Parse("2006-01-02T15-04", cfgDate)
		if err != nil {
			fmt.Println("Wrong time format provided.")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		nbEntries, err := wallabago.GetNumberOfTotalArticles()
		if err != nil {
			fmt.Println("Couldn't retrieve total number of articles.")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("There are", nbEntries, "saved on your wallabag instance.")

		url := wallabago.Config.WallabagURL +
			"/api/entries.json?perPage=" +
			strconv.Itoa(nbEntries) +
			"&detail=metadata" +
			"&sort=updated" +
			"&order=desc"
		response, err := wallabago.APICall(
			url,
			"GET",
			[]byte{},
		)
		if err != nil {
			fmt.Println("Couldn't retrieve articles from wallabag.")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var e wallabago.Entries
		err = json.Unmarshal(response, &e)
		if err != nil {
			fmt.Println("Bad format response from wallabag.")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		unread := 0
		if cfgUnread {
			unread = 1
		}
		starred := 0
		if cfgStarred {
			starred = 1
		}
		fmt.Println(
			"Will remove article older than",
			baseDate.Format("2006-01-02"),
			"and in status (unread, starred):",
			unread,
			starred,
		)
		var toRemove []wallabago.Item
		toRemove = ArticlesToRemove(e, baseDate, unread, starred, cfgForce)

		if len(toRemove) == 0 {
			fmt.Println("Nothing to delete, leaving")
			os.Exit(0)
		}
		PrintCandidateArticles(toRemove)

		if cfgDelete {
			DeleteArticles(toRemove)
		}
	},
}

func DeleteArticles(toRemove []wallabago.Item) {
	baseURL := wallabago.Config.WallabagURL + "/api/entries/"
	for i := 0; i < len(toRemove); i++ {
		url := baseURL + strconv.Itoa(toRemove[i].ID)
		response, err := wallabago.APICall(
			url,
			"DELETE",
			[]byte{},
		)
		if err != nil {
			fmt.Println("Couldn't delete entry", toRemove[i].ID)
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var item wallabago.Item
		err = json.Unmarshal(response, &item)
		if err != nil {
			fmt.Println("Bad format response from wallabag")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("Entry", item.Title, "(", item.URL, ") has been deleted.")
		time.Sleep(500 * time.Millisecond)
	}
}

func PrintCandidateArticles(toRemove []wallabago.Item) {
	fmt.Println("This command will remove", len(toRemove), "entries:")
	for i := 0; i < len(toRemove); i++ {
		status := "  "
		if toRemove[i].IsArchived == 0 {
			status = "🆕"
		}
		if toRemove[i].IsStarred == 1 {
			status += "⭐"
		} else {
			status += "  "
		}

		fmt.Println(
			"- ",
			toRemove[i].UpdatedAt.Format("2006-01-02"),
			status,
			toRemove[i].Title,
		)
	}
}

func ArticlesToRemove(
	e wallabago.Entries,
	baseDate time.Time,
	unread int,
	starred int,
	cfgForce bool, // Define cfgForce or pass it as a parameter to the function
) []wallabago.Item {
	var toRemove []wallabago.Item
	for i := 0; i < len(e.Embedded.Items); i++ {
		if e.Embedded.Items[i].UpdatedAt.Time.Before(baseDate) &&
			(cfgForce ||
				((e.Embedded.Items[i].IsArchived == 1 || unread == 1) &&
					(e.Embedded.Items[i].IsStarred == 0 || starred == 1))) {
			toRemove = append(toRemove, e.Embedded.Items[i])
		}
	}
	return toRemove
}
