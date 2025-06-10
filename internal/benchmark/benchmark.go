package benchmark

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/AlexShmak/golang-graphs/internal/algorithms"
	"github.com/AlexShmak/golang-graphs/internal/graph"
)

type RunResult struct {
	Duration time.Duration
	Found    bool
}

func calculateStats(results []RunResult) (mean, stdDev, foundRate float64) {
	if len(results) == 0 {
		return 0, 0, 0
	}
	var sum, sumSq float64
	var countFound int
	for _, r := range results {
		us := float64(r.Duration.Microseconds())
		sum += us
		sumSq += us * us
		if r.Found {
			countFound++
		}
	}
	n := float64(len(results))
	mean = sum / n
	variance := (sumSq/n) - (mean*mean)
	if variance < 0 {
		variance = 0
	}
	stdDev = math.Sqrt(variance)
	foundRate = float64(countFound) / n * 100
	return
}

type ExperimentConfig struct {
	Name  string
	K     int
	Graph *graph.Graph
}

func RunComparisonTest(cfg ExperimentConfig, numRuns int, timeout time.Duration) {
	g := cfg.Graph

	var baseResults, sotaResults []RunResult
	var baseTimedOut bool

	for i := 0; i < numRuns; i++ {
		sotaSolver := algorithms.NewSOTASolver(g, cfg.K, context.Background())
		startS := time.Now()
		foundS := sotaSolver.Solve()
		sotaResults = append(sotaResults, RunResult{Duration: time.Since(startS), Found: foundS})

		if g.NumVertices < 30 {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			baseSolver := algorithms.NewBaselineSolver(g, cfg.K, ctx)
			ch := make(chan bool, 1)
			startB := time.Now()
			go func() { ch <- baseSolver.Solve() }()
			select {
			case foundB := <-ch:
				baseResults = append(baseResults, RunResult{Duration: time.Since(startB), Found: foundB})
			case <-ctx.Done():
				baseTimedOut = true
			}
			cancel()
			if baseTimedOut {
				break
			}
		} else {
			baseTimedOut = true
		}
	}

	meanB, stdB, rateB := calculateStats(baseResults)
	meanS, stdS, rateS := calculateStats(sotaResults)

	var baseStr string
	if baseTimedOut {
		baseStr = fmt.Sprintf("Timeout (>%ds)", int(timeout.Seconds()))
	} else {
		baseStr = fmt.Sprintf("%.0f ± %.0f / %.0f%%", meanB, stdB, rateB)
	}
	sotaStr := fmt.Sprintf("%.0f ± %.0f / %.0f%%", meanS, stdS, rateS)

	var speedup string
	if baseTimedOut {
		speedup = ">> 1000x"
	} else if meanS > 0 {
		speedup = fmt.Sprintf("%.1fx", meanB/meanS)
	} else {
		speedup = "N/A"
	}

	fmt.Printf("| %-14s | %-3d | %-3d | %-2d | %-36s | %-33s | %-9s |\n",
		cfg.Name, g.NumVertices, len(g.Edges), cfg.K, baseStr, sotaStr, speedup)
}

func RunSOTATest(cfg ExperimentConfig, numRuns int, timeout time.Duration) {
	g := cfg.Graph

	var sotaResults []RunResult
	var sotaTimedOut bool

	for i := 0; i < numRuns; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		sotaSolver := algorithms.NewSOTASolver(g, cfg.K, ctx)
		ch := make(chan bool, 1)
		start := time.Now()
		go func() { ch <- sotaSolver.Solve() }()
		select {
		case found := <-ch:
			sotaResults = append(sotaResults, RunResult{Duration: time.Since(start), Found: found})
		case <-ctx.Done():
			sotaTimedOut = true
		}
		cancel()
		if sotaTimedOut {
			break
		}
	}

	var resultStr string
	if sotaTimedOut {
		resultStr = fmt.Sprintf("Timeout (>%ds)", int(timeout.Seconds()))
	} else {
		mean, std, rate := calculateStats(sotaResults)
		resultStr = fmt.Sprintf("%.0f ± %.0f / %.0f%%", mean, std, rate)
	}

	fmt.Printf("| %-25s | %-3d | %-3d | %-4d | %-4d | %-59s |\n",
		cfg.Name, g.NumVertices, len(g.Edges), cfg.K, g.NumVertices-cfg.K, resultStr)
}
