# Simulate Risk Limiting Audits

The Colorado Secretary of State (Jena Griswold) did an excellent
job getting ahead of rubbish conspiracy theories about voting
and election denial in 2020.

Part of the success of vote-by-mail in Colorado is that Colorado
does "Risk-limiting Audits" of ballots.
State of Colorado mails every eligible voter a paper ballot
that you mark by hand.
The ballot is both human readble and machine readable

State of Oregon has [similar procedures](https://sos.oregon.gov/elections/Pages/security.aspx),
as does [state of Washington](https://www.sos.wa.gov/elections/data-research/election-technology/post-election-audits).
Election denial and efforts to do in person, same day voting
are motivated by something other than desire for secure, reliable voting.

## A Gentle Introduction to Risk-Limiting Audits

The Colorado Secretary of State references a paper,
_A Gentle Introduction to Risk-Limiting Audits_,
by Mark Lindemann and Philip B. Stark,
from [IEEE SECURITY AND PRIVACY](https://ieeexplore.ieee.org/xpl/aboutJournal.jsp?punumber=8013),
Special Issue on Electronic Voting, 2012

The Coloraod web page links to a [preprint](https://www.stat.berkeley.edu/~stark/Preprints/gentle12.pdf)
of the paper, but as is sometimes the case,
the preprint differs significantly from the official, published version.
The preprint provokes confusion around the ballot auditing method,
in that is has some confusion about using percentages or proportions (i.e. 47% vs 0.47).

Get the official version if you can, matey!

## Build and Run the Simulation

It's in Go, it does not use any non-standard packages. It should be portable.

```
$ go build rla.go
```

That should leave you with an executable named `rla` if you're running a sane operating system.

### Options

* -b int
  * count of ballots (default 1000)
* -c string
  * Somethig like: `A:51,B:49` or `Smith:37,Wesson:20,Glock:43`.
  There's no preset limit on number of candidates or their names.
  The sum of the percentages of votes they got (37%, 20%, 43% in second example)
  does have to equal 100.0
* -f int
  * Voting fineness (default 1000). How many buckets to break up the probability distribution.
* -t float
  * tolerance for ballot audit, 0.0 chooses maximum tolerance
* -z 
  * Lie about who won. Instead of the winning candidate, returns the runner up,
  with the winning candidate's vote count

### Running it

```
$ ./rla -c A:51.0,B:49.0 -t .10 -b 100000
Desired results for 100000 ballots:
Candidate A: 51.00%
Candidate B: 49.00%
Generated results:
Candidate       Count   Desired Generated
A               50830   51.00   50.83
B               49170   49.00   49.17
Vote counting results:
A declared winner, with 50.8300%
Audit results:
17054 ballots of 100000 examined
Audit confirms winner
```

This simulates an election with two candidates, imaginatively named "A" and "B".
I wanted "A" to get about 51% of the votes, and "B" to get 49%
"A" got a simulated 50.83%, and "B" got 49.17%.

I told the program to use a tolerance `t` of 0.10%.
Because it was a fairly tight contest,
the ballot audit had to pull 17,054 ballots, or 17.054%,
to get to 90% sure that "A" really did win.

The same election with `-z` flag, which should tel you to do a hand recount:

```
$ ./rla -c A:51.0,B:49.0 -t .10 -b 100000  -z
Desired results for 100000 ballots:
Candidate A: 51.00%
Candidate B: 49.00%
Generated results:
Candidate       Count   Desired Generated
A               51019   51.00   51.02
B               48981   49.00   48.98
Vote counting results:
B declared winner, with 51.0190%
Audit results:
9140 ballots of 100000 examined
Hand recount to confirm
```

The audit does detect that "B" did not actually receive 51% of the votes,
and indicates a recount should occur.

## Experience

Futzing around with the `rla` program does show that the number of ballots
to check in an audit goes up when the winner gets closer to 50% of the votes.
The audit does reliably detect when the winner has more than 50% of the vote,
or when the runnerup is claimed as the winner.
