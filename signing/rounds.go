package keygen

import (
	"github.com/binance-chain/tss-lib/tss"
)

const (
	TaskName = "signing"
)

type (
	base struct {
		*tss.Parameters
		// save    *LocalPartySaveData
		// temp    *LocalPartyTempData
		out     chan<- tss.Message
		ok      []bool // `ok` tracks parties which have been verified by Update()
		started bool
		number  int
	}
	preparation struct {
		*base
	}
	round1 struct {
		*preparation
	}
	round2 struct {
		*round1
	}
	round3 struct {
		*round2
	}
	round4 struct {
		*round3
	}
	round5 struct {
		*round4
	}
	round6 struct {
		*round5
	}
	round7 struct {
		*round6
	}
	round8 struct {
		*round7
	}
	round9 struct {
		*round8
	}
	finalization struct {
		*round9
	}
)

var (
	_ tss.Round = (*preparation)(nil)
	_ tss.Round = (*round1)(nil)
	_ tss.Round = (*round2)(nil)
	_ tss.Round = (*round3)(nil)
	_ tss.Round = (*round4)(nil)
	_ tss.Round = (*round5)(nil)
	_ tss.Round = (*round6)(nil)
	_ tss.Round = (*round7)(nil)
	_ tss.Round = (*round8)(nil)
	_ tss.Round = (*round9)(nil)
	_ tss.Round = (*finalization{})(nil)
)

// ----- //

func (round *base) Params() *tss.Parameters {
	return round.Parameters
}

func (round *base) RoundNumber() int {
	return round.number
}

// CanProceed is inherited by other rounds
func (round *base) CanProceed() bool {
	if !round.started {
		return false
	}
	for _, ok := range round.ok {
		if !ok {
			return false
		}
	}
	return true
}

// WaitingFor is called by a Party for reporting back to the caller
func (round *base) WaitingFor() []*tss.PartyID {
	Ps := round.Parties().Parties()
	ids := make([]*tss.PartyID, 0, len(round.ok))
	for j, ok := range round.ok {
		if ok { continue }
		ids = append(ids, Ps[j])
	}
	return ids
}

func (round *base) WrapError(err error, culprits ...*tss.PartyID) *tss.Error {
	return tss.NewError(err, TaskName, round.number, round.PartyID(), culprits...)
}

// `ok` tracks parties which have been verified by Update()
func (round *base) resetOk() {
	for j := range round.ok {
		round.ok[j] = false
	}
}
