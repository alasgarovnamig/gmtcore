package request

import "gmtcore/search"

type SearchRequestDto struct {
	Criteria []search.Criteria
	Preloads []string
}

func (s *SearchRequestDto) MarkedDto() {}
