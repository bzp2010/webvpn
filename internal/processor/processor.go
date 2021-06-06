package processor

import (
	"net/http"

	"github.com/gogf/gf/container/glist"
)

type Processor interface {
	RequestProcess(*http.Request) bool
	ResponseProcess(*http.Response) bool
}

type Manager struct {
	list *glist.List
}

func NewManager() *Manager {
	return &Manager{
		glist.New(),
	}
}

func (m *Manager) Use(p ...Processor) {
	for _, v := range p {
		m.list.PushBack(v)
	}
}

func (m *Manager) DoRequestProcess(req *http.Request) {
	m.list.IteratorAsc(func(e *glist.Element) bool {
		return e.Value.(Processor).RequestProcess(req)
	})
}

func (m *Manager) DoResponseProcess(resp *http.Response) {
	m.list.IteratorDesc(func(e *glist.Element) bool {
		return e.Value.(Processor).ResponseProcess(resp)
	})
}
