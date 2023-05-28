package types

type AccountBalance struct {
	CreditsPosted  int
	CreditsPending int
	CreditsTotal   int

	DebitsPosted  int
	DebitsPending int
	DebitsTotal   int
}
