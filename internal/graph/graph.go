package graph

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Graph struct {
	Adj         map[int][]int
	Edges       [][2]int
	NumVertices int
}

func NewGraphFromEdges(edges [][2]int, numVertices int) *Graph {
	g := &Graph{
		Adj:         make(map[int][]int, numVertices),
		Edges:       edges,
		NumVertices: numVertices,
	}
	for i := 0; i < numVertices; i++ {
		g.Adj[i] = []int{}
	}
	for _, e := range edges {
		g.Adj[e[0]] = append(g.Adj[e[0]], e[1])
		g.Adj[e[1]] = append(g.Adj[e[1]], e[0])
	}
	return g
}

func GenerateGridGraph(rows, cols int) *Graph {
	n := rows * cols
	edges := [][2]int{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			i := r*cols + c
			if r+1 < rows {
				edges = append(edges, [2]int{i, (r + 1) * cols + c})
			}
			if c+1 < cols {
				edges = append(edges, [2]int{i, r*cols + (c + 1)})
			}
		}
	}
	return NewGraphFromEdges(edges, n)
}

func LoadTxtGraph(path string) *Graph {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Не удалось открыть файл %s: %v", path, err)
	}
	defer file.Close()
	edges := [][2]int{}
	nodeMap := make(map[int]int)
	maxNodeID := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		u, _ := strconv.Atoi(parts[0])
		v, _ := strconv.Atoi(parts[1])
		if _, ok := nodeMap[u]; !ok {
			nodeMap[u] = maxNodeID
			maxNodeID++
		}
		if _, ok := nodeMap[v]; !ok {
			nodeMap[v] = maxNodeID
			maxNodeID++
		}
		edges = append(edges, [2]int{nodeMap[u], nodeMap[v]})
	}
	return NewGraphFromEdges(edges, maxNodeID)
}

func LoadGMLGraph(path string) *Graph {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Не удалось открыть файл %s: %v", path, err)
	}
	defer file.Close()
	bytes, _ := os.ReadFile(path)
	content := string(bytes)
	nodeRe := regexp.MustCompile(`\s*id\s+(\d+)`)
	sourceRe := regexp.MustCompile(`\s*source\s+(\d+)`)
	targetRe := regexp.MustCompile(`\s*target\s+(\d+)`)
	nodeMatches := nodeRe.FindAllStringSubmatch(content, -1)
	numVertices := 0
	idMap := make(map[int]int)
	for _, match := range nodeMatches {
		id, _ := strconv.Atoi(match[1])
		if _, ok := idMap[id]; !ok {
			idMap[id] = numVertices
			numVertices++
		}
	}
	edgeSections := strings.Split(content, "edge")
	edges := [][2]int{}
	for _, section := range edgeSections[1:] {
		sourceMatch := sourceRe.FindStringSubmatch(section)
		targetMatch := targetRe.FindStringSubmatch(section)
		if len(sourceMatch) > 1 && len(targetMatch) > 1 {
			uID, _ := strconv.Atoi(sourceMatch[1])
			vID, _ := strconv.Atoi(targetMatch[1])
			u := idMap[uID]
			v := idMap[vID]
			edges = append(edges, [2]int{u, v})
		}
	}
	return NewGraphFromEdges(edges, numVertices)
}
