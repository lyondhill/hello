package main

import (
	"github.com/mmcdole/gofeed"
	"fmt"
	// "encoding/json"
	// "os"
	// "log"
	// "bytes"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://feeds.reuters.com/reuters/businessNews")
	if err != nil {
		panic(err)
	}
	// b, err := json.Marshal(feed)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var out bytes.Buffer
	// json.Indent(&out, b, "=", "\t")
	// out.WriteTo(os.Stdout)
	for _, item := range feed.Items {
		fmt.Println("item.Title", item.Title)
		fmt.Println("item.Description", item.Description)
		fmt.Println("item.Content", item.Content)
		fmt.Println("item.Link", item.Link)
		fmt.Println("item.Updated", item.Updated)
		fmt.Println("item.Published", item.Published)
		fmt.Println("item.Author", item.Author)
		fmt.Println("item.GUID", item.GUID)
	}
}
