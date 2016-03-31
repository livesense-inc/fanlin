package imgproxyerr

const (
	ERROR   = "error"
	WARNING = "warning"
)

// やばそうやつから配列の前に積んでいく
var level = [...]string{
	ERROR,
	WARNING,
}

func getLevel(Type string) int {
	for l, t := range level {
		if t == Type {
			return l
		}
	}
	// 存在しないのが来た時わからないので危険側に倒す
	// 仕様上からの見過ごし防止の為
	return -1
}

func New(t string, err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Err); ok {
		if e.cmp(t) {
			return &Err{e.Type, e}
		}
	}
	return &Err{t, err}
}

type Err struct {
	Type string
	Err  error
}

func (e *Err) Error() string {
	return e.Err.Error()
}

func (e *Err) cmp(v string) bool {
	return getLevel(e.Type)-getLevel(v) <= 0
}
