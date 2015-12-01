package browscap_go

import (
	"bytes"
	"github.com/hamaxx/browscap_go/re0"
)

type expression struct {
	Name string
	exp  re0.Expression
}

func newRegexpExpression(val string) *expression {
	return &expression{
		Name: val,
		exp:  re0.Compile(bytes.ToLower([]byte(val))),
	}
}

func (self *expression) Match(val []byte) bool {
	return self.exp.Match(val)
}

type expressionByNameLen []*expression

func (el expressionByNameLen) Len() int {
	return len(el)
}

func (el expressionByNameLen) Less(i, j int) bool {
	return len(el[i].Name) > len(el[j].Name)
}

func (el expressionByNameLen) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}
