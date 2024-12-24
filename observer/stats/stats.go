package stats

type Kind int

const (
	KindTotalConns   Kind = 1
	KindCurrentConns Kind = 2
	KindInputBytes   Kind = 3
	KindOutputBytes  Kind = 4
	KindTotalErrs    Kind = 5
)

type Stats interface {
	Add(kind Kind, n int64)
	Get(kind Kind) uint64
	IsUpdated() bool
	Reset()
}
