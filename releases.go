package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/html"
)

func releases(ctx context.Context, cmd *cli.Command) error {
	url := gvmhost + "/dl/"
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return cli.Exit(fmt.Sprintf("Could not fetch Go releases: %d", res.StatusCode), 1)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}

	nodes := doc.Find(".toggleVisible").Nodes
	nodes = append(nodes, doc.Find(".toggle").Nodes...)
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}

	for _, node := range nodes {
		selector := &goquery.Selection{
			Nodes: []*html.Node{node},
		}
		if version, ok := selector.Attr("id"); ok {
			if !strings.HasPrefix(version, "go") {
				continue
			}

			fmt.Println(strings.TrimPrefix(version, "go"))
		}
	}

	return nil
}
