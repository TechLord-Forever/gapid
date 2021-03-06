// Copyright (C) 2017 Google Inc.
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

package gles

import (
	"strings"

	"github.com/google/gapid/gapis/vertex"
)

var semanticPatterns = []struct {
	pattern  string
	semantic vertex.Semantic_Type
}{
	// Ordered from highest priority to lowest
	{"position", vertex.Semantic_Position},
	{"normal", vertex.Semantic_Normal},
	{"tangent", vertex.Semantic_Tangent},
	{"bitangent", vertex.Semantic_Bitangent},
	{"binormal", vertex.Semantic_Bitangent},
	{"texcoord", vertex.Semantic_Texcoord},
	{"pos", vertex.Semantic_Position},
	{"uv", vertex.Semantic_Texcoord},
	{"vertex", vertex.Semantic_Position},
}

func guessSemantics(vb *vertex.Buffer) {
	taken := map[vertex.Semantic_Type]bool{}
	for _, p := range semanticPatterns {
		if taken[p.semantic] {
			continue
		}
		for _, s := range vb.Streams {
			if strings.Contains(strings.ToLower(s.Name), p.pattern) {
				s.Semantic.Type = p.semantic
				taken[p.semantic] = true
				break
			}
		}
	}
}
