package schema

import (
	"strings"
	"time"
)

const (
	maxPageSize     = 1000
	defaultPageSize = 10
)

type Pageable struct {
	Size   *int `query:"size"   oai:"description=返回数据数量;default=10;maximum=1000"`
	Offset *int `query:"offset" oai:"description=数据偏移量"`
}

func (p *Pageable) GetLimit() int {
	if p.Size == nil {
		return defaultPageSize
	}
	if *p.Size >= maxPageSize {
		return maxPageSize
	}
	return *p.Size
}

func (p *Pageable) GetOffset() int {
	if p.Offset == nil {
		return 0
	}
	return *p.Offset
}

type Sortable struct {
	OrderBy *string `query:"order_by" oai:"description=排序字段, 如: id desc/asc"`
}

func (s *Sortable) GetOrderField() (field string, asc bool) {
	if s.OrderBy == nil {
		return "", true
	}
	parts := strings.Split(*s.OrderBy, " ")
	if len(parts) == 1 {
		return parts[1], true
	}
	if len(parts) == 2 {
		switch parts[1] {
		case "asc", "ascend", "ascending":
			return parts[0], true
		case "desc", "descend", "descending":
			return parts[0], false
		}
	}
	return "", true
}

type SortableBody struct {
	OrderBy *string `json:"order_by" oai:"description=排序字段, 如: id desc/asc"`
}

func (s *SortableBody) GetOrderField() (field string, asc bool) {
	if s.OrderBy == nil {
		return "", true
	}
	parts := strings.Split(*s.OrderBy, " ")
	if len(parts) == 1 {
		return parts[1], true
	}
	if len(parts) == 2 {
		switch parts[1] {
		case "asc", "ascend", "ascending":
			return parts[0], true
		case "desc", "descend", "descending":
			return parts[0], false
		}
	}
	return "", true
}

type (
	TimeSeries     []TimeSeriesItem
	TimeSeriesItem struct {
		Time time.Time `json:"time" oai:"description=时间"`
		V    int       `json:"v"    oai:"description=值"`
	}
)

type PageableBody struct {
	Size   *int `json:"size"   oai:"description=返回数据数量;default=10;maximum=1000"`
	Offset *int `json:"offset" oai:"description=数据偏移量"`
}

func (p *PageableBody) GetLimit() int {
	if p.Size == nil {
		return defaultPageSize
	}
	if *p.Size >= maxPageSize {
		return maxPageSize
	}
	return *p.Size
}

func (p *PageableBody) GetOffset() int {
	if p.Offset == nil {
		return 0
	}
	return *p.Offset
}
