//  Copyright 2013 Walter Schulze
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

package compose

import (
	"fmt"
	"github.com/awalterschulze/katydid/funcs"
	"reflect"
)

type errNotConst struct {
	f     interface{}
	field reflect.Value
}

func (this *errNotConst) Error() string {
	return fmt.Sprintf("%s has constant %s which has a variable parameter", reflect.ValueOf(this.f).Elem().Type().Name(), this.field.Elem().Type())
}

func TrimBool(f funcs.Bool) (funcs.Bool, error) {
	trimmed, err := trim(f)
	if err != nil {
		return nil, err
	}
	return trimmed.(funcs.Bool), nil
}

func trim(f interface{}) (interface{}, error) {
	if reflect.TypeOf(f).Implements(varTyp) {
		return f, nil
	}
	if reflect.TypeOf(f).Implements(constTyp) {
		return f, nil
	}
	this := reflect.ValueOf(f).Elem()
	trimable := true
	for i := 0; i < this.NumField(); i++ {
		if _, ok := this.Field(i).Type().MethodByName("Eval"); !ok {
			continue
		}
		if this.Field(i).Elem().Type().Implements(varTyp) {
			trimable = false
			continue
		}
		trimmed, err := trim(this.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		this.Field(i).Set(reflect.ValueOf(trimmed))
		if !this.Field(i).Elem().Type().Implements(constTyp) {
			if funcs.IsConst(this.Field(i).Type()) {
				return nil, &errNotConst{f, this.Field(i)}
			}
			trimable = false
		}
	}
	if !trimable {
		return f, nil
	}
	if inits, ok := f.(funcs.Init); ok {
		err := inits.Init()
		if err != nil {
			return nil, err
		}
	}
	return funcs.NewConst(reflect.ValueOf(f).MethodByName("Eval").Call(nil)[0].Interface()), nil
}
