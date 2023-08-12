package util

import (
	"fmt"
	"github.com/dominikbraun/graph"
)

func BFS[K comparable, T any](g graph.Graph[K, T], start K, forward bool, visit func(K, K) bool) error {
	var (
		adjacencyMap map[K]map[K]graph.Edge[K]
		err          error
	)
	if forward {
		adjacencyMap, err = g.AdjacencyMap()
		if err != nil {
			return fmt.Errorf("could not get adjacency map: %w", err)
		}
	} else {
		adjacencyMap, err = g.PredecessorMap()
		if err != nil {
			return fmt.Errorf("could not get adjacency map: %w", err)
		}
	}

	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[start] = true
	queue = append(queue, start)

	for len(queue) > 0 {
		currentHash := queue[0]

		queue = queue[1:]

		for adjacency := range adjacencyMap[currentHash] {
			// Stop traversing the graph if the visit function returns true.
			if stop := visit(currentHash, adjacency); stop {
				break
			}

			if _, ok := visited[adjacency]; !ok {
				visited[adjacency] = true
				queue = append(queue, adjacency)
			}
		}

	}

	return nil
}
