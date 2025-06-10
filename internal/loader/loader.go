package loader

import (
	"github.com/AlexShmak/golang-graphs/internal/graph"
)

func LoadTXT(path string) *graph.Graph {
	return graph.LoadTxtGraph(path)
}

func LoadGML(path string) *graph.Graph {
	return graph.LoadGMLGraph(path)
}
