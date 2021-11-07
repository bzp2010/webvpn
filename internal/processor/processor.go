/*
 * Copyright (C) 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
