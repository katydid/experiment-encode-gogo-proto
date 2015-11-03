//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package debug

import (
	"fmt"
	"github.com/katydid/katydid/serialize"
	"io"
	"math/rand"
	"strings"
	"time"
)

func getValue(parser serialize.Parser) interface{} {
	var v interface{}
	var err error
	v, err = parser.Int()
	if err == nil {
		return v
	}
	v, err = parser.Uint()
	if err == nil {
		return v
	}
	v, err = parser.Double()
	if err == nil {
		return v
	}
	v, err = parser.Bool()
	if err == nil {
		return v
	}
	v, err = parser.String()
	if err == nil {
		return v
	}
	v, err = parser.Bytes()
	if err == nil {
		return v
	}
	return nil
}

type Node struct {
	Label    string
	Children Nodes
}

func (this Node) String() string {
	if len(this.Children) == 0 {
		return this.Label
	}
	return this.Label + ":" + this.Children.String()
}

func (this Node) Equal(that Node) bool {
	if this.Label != that.Label {
		return false
	}
	if !this.Children.Equal(that.Children) {
		return false
	}
	return true
}

type Nodes []Node

func (this Nodes) String() string {
	ss := make([]string, len(this))
	for i := range this {
		ss[i] = this[i].String()
	}
	return "{" + strings.Join(ss, ",") + "}"
}

func (this Nodes) Equal(that Nodes) bool {
	if len(this) != len(that) {
		return false
	}
	for i := range this {
		if !this[i].Equal(that[i]) {
			return false
		}
	}
	return true
}

func Walk(parser serialize.Parser) Nodes {
	a := make(Nodes, 0)
	for {
		if err := parser.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		value := getValue(parser)
		if parser.IsLeaf() {
			a = append(a, Node{fmt.Sprintf("%v", value), nil})
		} else {
			name := fmt.Sprintf("%v", value)
			parser.Down()
			v := Walk(parser)
			parser.Up()
			a = append(a, Node{name, v})
		}
	}
	return a
}

func NewRand() Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Rand interface {
	Intn(n int) int
}

func RandomWalk(parser serialize.Parser, r Rand, next, down int) Nodes {
	a := make(Nodes, 0)
	for {
		if r.Intn(next) == 0 {
			break
		}
		if err := parser.Next(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		value := getValue(parser)
		if parser.IsLeaf() {
			a = append(a, Node{fmt.Sprintf("%#v", value), nil})
		} else {
			name := fmt.Sprintf("%#v", value)
			var v Nodes
			if r.Intn(down) != 0 {
				parser.Down()
				v = RandomWalk(parser, r, next, down)
				parser.Up()
			}
			a = append(a, Node{name, v})
		}
	}
	return a
}
