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
	minimizeTolerance := flag.Bool("m", false, "minimize tolerance in audit")
	flag.Parse()

	candidates, err := parseCandidates(*candidateProportions)
	if err != nil {
		log.Fatal(err)
	}

	for _, candidate := range candidates {
		fmt.Printf("Candidate %s: %.02f%%\n", candidate.name, candidate.percentage)
	}

	voting := createVoting(candidates, *votingFineness)

	ballots := createVotes(voting, *ballotCount, *votingFineness)

	winner, winningPercentage := countVotes(candidates, *ballotCount, *lie)
	fmt.Printf("%s declared winner, %02f%%\n", winner, winningPercentage)

	recount := auditBallots(ballots, *ballotCount, winner, winningPercentage, *minimizeTolerance)

	if recount {
		fmt.Printf("hand recount to confirm\n")
	}
}

// auditVote runs the Ballot-polling audit of section III.a of
// A Gentle Introduction to Risk-limiting Audits,
// Mark Lindemann, Philip B. Stark,
// IEEE SECURITY AND PRIVACY, SPECIAL ISSUE ON ELECTRONIC VOTING, 2012. LAST EDITED 16 MARCH 2012.
func auditBallots(ballots []string, ballotCount int, winner string, winningPercentage float64, minimizeTolerance bool) bool {

	// tolerance - "Let t be a positive "tolerance" small enough
	// that when t is subtracted from the winner's vote share
	// s, the difference is still greater than 50%.
	t := winningPercentage - 49.9
	if minimizeTolerance {
		// This is harder
	}

	matchFactor := (winningPercentage - t) / 50.0
	notMatchFactor := 1.0 - matchFactor
	fmt.Printf("Audit, tolerance %.04f%%, match %.05f, not match %.05f\n", t, matchFactor, notMatchFactor)

	T := 1.0

	for ballotsExamined := 1; true; ballotsExamined++ {
		// "2) Select a ballot at random from the ballots cast in the contest"
		// Lindemann & Stark suggest a programmatic random number genertor,
		// I'm going to use a built-in. We don't have to track whether we've seen
		// the ballot before.
		ballotNumber := rand.Intn(ballotCount)

		// Steps (4) and (5), multiply T by a factor
		factor := matchFactor
		if ballots[ballotNumber] != winner {
			factor = notMatchFactor
		}
		T *= factor

		fmt.Printf("audit ballot %d: choice %q, factor %.05f, T now %.05f\n", ballotsExamined, ballots[ballotNumber], factor, T)

		// "6) If T > 9.9, the audit has provided strong evidence that
		// the reported outcome is correct: Stop."
		if T > 9.9 {
			break // to end-of-function return, keep compiler happy
		}

		// "7) If T < 0.011, perform a full hand count to determine
		// who won. Otherwise, return to step 2."
		if T < 0.011 {
			return true
		}
	}

	return false
}

func parseCandidates(candidateProportions string) ([]*candidate, error) {
	if candidateProportions == "" {
		return nil, errors.New("no candidates")
	}

	fields := strings.Split(candidateProportions, ",")

	candidates := make([]*candidate, 0)

	var sum float64
	over50Percent := false

	for _, f := range fields {
		g := strings.Split(f, ":")
		if p, err := strconv.ParseFloat(g[1], 64); err == nil {
			c := &candidate{
				name:       g[0],
				percentage: p,
			}
			candidates = append(candidates, c)
			sum += p
			if p > 50.00 {
				over50Percent = true
			}
		} else {
			log.Printf("# Candidate %q: %v\n", f, err)
		}
	}

	if !over50Percent {
		return nil, errors.New("no candidate received over 50%% of vote\n")
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

func countVotes(candidates []*candidate, ballotCount int, lie bool) (string, float64) {

	winner := candidates[0].name
	votes := candidates[0].votes
	winningPercentage := float64(votes) / float64(ballotCount) * 100.

	runnerUp := candidates[1].name

	for _, c := range candidates {
		if c.votes > votes {
			runnerUp = winner
			winner = c.name
			votes = c.votes
			winningPercentage = float64(c.votes) / float64(ballotCount) * 100.

		}
		fmt.Printf("%s\t%d\t%.02f\t%.02f\n",
			c.name,
			c.votes,
			c.percentage,
			float64(c.votes)/float64(ballotCount)*100.,
		)
	}

	if lie {
		return runnerUp, winningPercentage
	}

	return winner, winningPercentage
}
