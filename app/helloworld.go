package main

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"cloud.google.com/go/bigquery"
	"os"
)

var proj string

func main() {

	http.HandleFunc("/", indexHandler)

	proj = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		proj = "chainlink-marketing-roi"
	}

	//GAE env sets port env - dynamically allocated.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	num := getNumVisitors()

	fmt.Fprint(w, "There were %d visitors",num)
}

type Visitors struct {
        Visitors int64 `bigquery:"visitors"`
}

func getNumVisitors() int64{
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		log.Fatal(err)
        os.Exit(1)
    }
    query := client.Query(
    	`SELECT
    	COUNT(1) as visitors
        FROM ` + "`bigquery-public-data.stackoverflow.posts_questions`" + `;`)
     iter, read_error := query.Read(ctx)
     if read_error != nil {
     	log.Fatal(err)
     }
     for {
     	var row Visitors
        iterError := iter.Next(&row)
        if iterError == iterator.Done {
        	return 0
        }
        if iterError != nil {
        	log.Fatal(iterError)
        	return -1
        }
        fmt.Print("unique visitors: %d\n",row.Visitors)
        return row.Visitors
    }
}