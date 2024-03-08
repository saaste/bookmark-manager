package bookmarks

import "strings"

type ConditionBuilder struct {
	conditions []string
}

func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{
		conditions: make([]string, 0),
	}
}

func (cb *ConditionBuilder) Add(condition string) {
	cb.conditions = append(cb.conditions, condition)
}

func (cb *ConditionBuilder) String() string {
	return "WHERE " + strings.Join(cb.conditions, " AND ")
}
