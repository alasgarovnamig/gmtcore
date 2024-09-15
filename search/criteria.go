package search

type Operation int

const (
	GreaterThan Operation = iota + 1
	LessThan
	GreaterThanEqual
	LessThanEqual
	NotEqual
	Equal
	In
	NotIn
	Match
	MatchStart
	MatchEnd
	JoinChild      // Çocuk tablo ile join
	JoinGrandChild // Torun tablo ile join
	AnyOf          // Bir koleksiyon içindeki herhangi bir değerle eşleşme
	IsMember       // Bir alanın koleksiyon üyesi olup olmadığını kontrol eder
)

type Criteria struct {
	Key           string
	ChildKey      string
	GrandChildKey string
	Value         interface{}
	Operation     Operation
}
