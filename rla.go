package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type candidate struct {
	name       string
	percentage float64
	votes      int
}

func main() {
	candidateProportions := flag.String("c", "", "x:%,y:%(,z:%)")
	lie := flag.Bool("z", false, "maybe lie about who won")
	ballotCount := flag.Int("b", 1000, "count of ballots")
	votingFineness := flag.Int("f", 1000, "voting fineness")
	flag.Parse()

	candidates, err := parseCandidates(*candidateProportions)
	if err != nil {
		log.Fatal(err)
	}

	for _, candidate := range candidates {
		fmt.Printf("Candidate %s: %.02f%%\n", candidate.name, candidate.percentage)
	}

	voting := createVoting(candidates, *votingFineness)

	_ = createVotes(voting, *ballotCount, *votingFineness)

	winner := countVotes(candidates, *ballotCount, *lie)
	fmt.Printf("%s declared winner\n", winner)

	// auditVotes(votes, *ballotCount)
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
				percentage: p,
			}
			candidates = append(candidates, c)
			sum += p
		} else {
			log.Printf("# Candidate %q: %v\n", f, err)
		}
	}

	if sum < 99. || sum > 100. {
		return nil, fmt.Errorf("sum of percentages %.02f, != 100.\n", sum)
	}

	return candidates, nil
}

func createVoting(candidates []*candidate, fineness int) []*candidate {

	votes := make([]*candidate, fineness)
	i := 0

	factor := float64(fineness) / 100.0

	for _, c := range candidates {
		x := int(c.percentage * factor)
		for ; x > 0; x-- {
			votes[i] = c
			i++
		}
	}

	return votes
}

func createVotes(voting []*candidate, ballotCount int, fineness int) []string {

	votes := make([]string, ballotCount)

	for i := 0; i < ballotCount; i++ {
		j := rand.Intn(fineness)
		voting[j].votes++
		votes[i] = voting[j].name
	}

	return votes
}

func countVotes(candidates []*candidate, ballotCount int, lie bool) string {

	winner := candidates[0].name
	votes := candidates[0].votes
	for _, c := range candidates {
		if c.votes > votes {
			winner = c.name
			votes = c.votes
		}
		fmt.Printf("%s\t%d\t%.02f\t%.02f\n", c.name, c.votes, c.percentage, float64(c.votes)/float64(ballotCount)*100.)
	}
	return winner
}
