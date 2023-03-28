// Package cytoscape visualize resource graph
package cytoscape

import (
	"bytes"
	"embed"
	"html/template"
	"teredix/pkg/storage"
)

//go:embed index.html
var indexHTML embed.FS

// PageData hold page data
type PageData struct {
	Elements  Elements
	InlineCSS template.HTMLAttr
}

// Elements holds elements
type Elements struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// NodeData holds node data
type NodeData struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Node holds node
type Node struct {
	Data NodeData `json:"data"`
}

// EdgeData holds edge data
type EdgeData struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// Edge hold edge information
type Edge struct {
	Data EdgeData `json:"data"`
}

// Cytoscape visualization library
type Cytoscape struct {
	storage storage.Storage
}

// NewCytoscapa construct cytoscape
func NewCytoscapa(st storage.Storage) *Cytoscape {
	return &Cytoscape{storage: st}
}

// Display generate resource graph and return html page
func (c *Cytoscape) Display() (string, error) {

	resources, err := c.storage.GetResources()
	if err != nil {
		return "", err
	}

	var nodes []Node

	for _, r := range resources {
		nodes = append(nodes, Node{Data: NodeData{
			ID:    r.UUID,
			Label: r.Name,
		}})
	}

	// edges

	relations, err := c.storage.GetRelations()
	if err != nil {
		return "", err
	}

	var edges []Edge
	for _, e := range relations {
		for k, v := range e {
			edges = append(edges, Edge{
				Data: EdgeData{
					Source: k,
					Target: v,
				},
			})
		}
	}

	elements := Elements{
		Nodes: nodes,
		Edges: edges,
	}

	tmpl, err := template.ParseFS(indexHTML, "index.html")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, PageData{Elements: elements, InlineCSS: `width: 100%; height: 100%; position: absolute; top: 0; left: 0;`})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
