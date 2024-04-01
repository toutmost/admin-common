// Copyright 2023 The Ryan SU Authors (https://github.com/suyuan32). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// StateMixin is the mixin with state field which
// is bool type.
type StateMixin struct {
	mixin.Schema
}

func (StateMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("state").
			Default(true).
			Optional().
			Comment("State true: normal false: ban | 状态 true 正常 false 禁用").
			Annotations(entsql.WithComments(true)),
	}
}
