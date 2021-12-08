// Copyright © SurrealDB Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cork

import (
	"reflect"
	"sync"
)

const (
	maybe int = iota
	yes
	no
)

var c cache

type cache struct {
	cl sync.RWMutex
	c  map[reflect.Type]int
	sl sync.RWMutex
	s  map[reflect.Type]int
	ml sync.RWMutex
	m  map[reflect.Type][]*field
}

func init() {
	c = cache{
		c: make(map[reflect.Type]int),
		s: make(map[reflect.Type]int),
		m: make(map[reflect.Type][]*field),
	}
}

func (c *cache) Has(t reflect.Type) bool {
	c.ml.RLock()
	val := c.m[t] != nil
	c.ml.RUnlock()
	return val
}

func (c *cache) Get(t reflect.Type) []*field {
	c.ml.RLock()
	val := c.m[t]
	c.ml.RUnlock()
	return val
}

func (c *cache) Set(t reflect.Type, fls []*field) {
	c.ml.Lock()
	c.m[t] = fls
	c.ml.Unlock()
}

func (c *cache) Corkable(t reflect.Type) bool {
	c.cl.RLock()
	switch c.c[t] {
	case yes:
		c.cl.RUnlock()
		return true
	case no:
		c.cl.RUnlock()
		return false
	default:
		c.cl.RUnlock()
		c.cl.Lock()
		defer c.cl.Unlock()
		if t.Implements(typeCorker) {
			c.c[t] = yes
			return true
		} else {
			c.c[t] = no
			return false
		}
	}
}

func (c *cache) Selfable(t reflect.Type) bool {
	c.sl.RLock()
	switch c.s[t] {
	case yes:
		c.sl.RUnlock()
		return true
	case no:
		c.sl.RUnlock()
		return false
	default:
		c.sl.RUnlock()
		c.sl.Lock()
		defer c.sl.Unlock()
		if t.Implements(typeSelfer) {
			c.s[t] = yes
			return true
		} else {
			c.s[t] = no
			return false
		}
	}
}
