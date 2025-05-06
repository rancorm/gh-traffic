package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/cli/go-gh/v2/pkg/api"
)

// Create a combined struct for JSON output
type TrafficData struct {
	Repository string `json:"repository"`
	Views struct {
		Count   int `json:"count"`
		Uniques int `json:"uniques"`
	} `json:"views"`
	Clones struct {
		Count   int `json:"count"`
		Uniques int `json:"uniques"`
	} `json:"clones"`
}

func main() {
	owner := flag.String("owner", "", "Repository owner")
	repo := flag.String("repo", "", "Repository name")
	outputJSON := flag.Bool("json", false, "Output as JSON")
	
	flag.Parse()

	var repoOwner, repoName string

	if *owner != "" && *repo != "" {
		repoOwner = *owner
		repoName = *repo
	} else {
		repoCtx, err := repository.Current()
			
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not determine repository context:", err)
			os.Exit(1)
		}
		
		repoOwner = repoCtx.Owner
		repoName = repoCtx.Name
	}

	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating REST client: %v\n", err)
		os.Exit(1)
	}

	// Fetch views
	var views struct {
		Count     int `json:"count"`
		Uniques   int `json:"uniques"`
	}
	
	err = client.Get(fmt.Sprintf("repos/%s/%s/traffic/views", repoOwner, repoName), &views)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching views: %v\n", err)
		os.Exit(1)
	}

	// Fetch clones
	var clones struct {
		Count     int `json:"count"`
		Uniques   int `json:"uniques"`
	}
	
	err = client.Get(fmt.Sprintf("repos/%s/%s/traffic/clones", repoOwner, repoName), &clones)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching clones: %v\n", err)
		os.Exit(1)
	}
	
	// After fetching repo details, build traffic data struct.
	data := TrafficData{
		Repository: fmt.Sprintf("%s/%s", repoOwner, repoName),
		Views:      views,
		Clones:     clones,
	}

	// Output formatting
	if *outputJSON {
		json.NewEncoder(os.Stdout).Encode(data)
	} else {
		fmt.Printf("Repository: %s\n", data.Repository)
		fmt.Printf("    Clones: %d total, %d unique\n", data.Clones.Count, data.Clones.Uniques)
		fmt.Printf("  Visitors: %d total, %d unique\n", data.Views.Count, data.Views.Uniques)
	}
}
