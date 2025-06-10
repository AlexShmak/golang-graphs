package main

import (
	"fmt"
	"time"

	"github.com/AlexShmak/golang-graphs/internal/benchmark"
	"github.com/AlexShmak/golang-graphs/internal/graph"
	"github.com/AlexShmak/golang-graphs/internal/loader"
)

func main() {
	const numRuns = 3
	const baselineTimeout = 15 * time.Second

	fmt.Println("========= Эксперимент 1: Сравнение производительности Baseline и SOTA (N_RUNS=3) =========")
	tests1 := []benchmark.ExperimentConfig{
		{Name: "Решетка 3x3", Graph: graph.GenerateGridGraph(3, 3)},
		{Name: "Решетка 4x4", Graph: graph.GenerateGridGraph(4, 4)},
		{Name: "Решетка 5x5", Graph: graph.GenerateGridGraph(5, 5)},
		{Name: "Karate Club", Graph: loader.LoadTXT("data/karate.gml")},
	}
	fmt.Println("| Название теста | N   | E   | k  |              Время (Baseline, µs) |                  Время (SOTA, µs) | Ускорение |")
	fmt.Println("|----------------|-----|-----|----|--------------------------------------|-----------------------------------|-----------|")
	for _, cfg := range tests1 {
		cfg.K = cfg.Graph.NumVertices / 2
		benchmark.RunComparisonTest(cfg, numRuns, baselineTimeout)
	}

	fmt.Println("\n\n========= Эксперимент 2: Демонстрация масштабируемости SOTA (N_RUNS=3) =========")
	tests2 := []benchmark.ExperimentConfig{
		{Name: "Dolphins", Graph: loader.LoadGML("data/dolphins.gml"), K: 31},
		{Name: "Dolphins (k → N)", Graph: loader.LoadGML("data/dolphins.gml"), K: 55},
		{Name: "Football", Graph: loader.LoadGML("data/football.gml"), K: 57},
		{Name: "Football (k → N)", Graph: loader.LoadGML("data/football.gml"), K: 100},
	}
	const sotaTimeout = 30 * time.Second
	fmt.Println("| Название теста            | N   | E   | k    | N-k  |       Время (SOTA, µs) (среднее ± стд. откл. / % найдено) |")
	fmt.Println("|---------------------------|-----|-----|------|------|-------------------------------------------------------------|")
	for _, cfg := range tests2 {
		benchmark.RunSOTATest(cfg, numRuns, sotaTimeout)
	}
}
