package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type candidate struct {
	name       string
	proportion float64
}

func main() {
	candidateProportions := flag.String("c", "", "x:%,y:%(,z:%)")
	lie := flag.Bool("z", false, "maybe lie about who won")
	flag.Parse()

	candidates, err := parseCandidates(*candidateProportions)
	if err != nil {
		log.Fatal(err)
	}

	for _, candidate := range candidates {
		fmt.Printf("Candidate %s: %.02f%%\n", candidate.name, candidate.proportion)
	}

	if *lie {
		fmt.Printf("might lie about who won\n")
	}
}

func parseCandidates(candidateProportions string) ([]*candidate, error) {
	if candidateProportions == "" {
		return nil, errors.New("no candidates")
	}

	fields := strings.Split(candidateProportions, ",")

	candidates := make([]*candidate, 0)

	var sum float64

	for _, f := range fields {
		g := strings.Split(f, ":")
		if p, err := strconv.ParseFloat(g[1], 64); err == nil {
			c := &candidate{
				name:       g[0],
				proportion: p,
			}
			candidates = append(candidates, c)
			sum += p
		} else {
			log.Printf("Candidate %q: %v\n", f, err)
		}
	}

	if sum < 99. || sum > 100. {
		return nil, fmt.Errorf("sum of percentages %.02f, != 100.\n", sum)
	}

	return candidates, nil
}
