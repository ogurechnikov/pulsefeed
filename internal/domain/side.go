package domain

type Side int

const (
	SideUnknown Side = iota
	SideBuy
	SideSell
)

func (s Side) IsBuy() bool {
	return s == SideBuy
}

func (s Side) IsSell() bool {
	return s == SideSell
}

func (s Side) String() string {
	switch s {
	case SideBuy:
		return "BUY"
	case SideSell:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}
