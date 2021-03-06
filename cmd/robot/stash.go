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

package main

import (
	"flag"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/google/gapid/core/app"
	"github.com/google/gapid/core/data/search/script"
	"github.com/google/gapid/core/data/stash"
	stashgrpc "github.com/google/gapid/core/data/stash/grpc"
	"github.com/google/gapid/core/fault/cause"
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/net/grpcutil"
	"google.golang.org/grpc"
)

func init() {
	stashUpload := &app.Verb{
		Name:       "stash",
		ShortHelp:  "Upload a file to the stash",
		ShortUsage: "<filenames>",
		Run:        doUpload(stashUploader{}),
	}
	uploadVerb.Add(stashUpload)
	stashSearch := &app.Verb{
		Name:       "stash",
		ShortHelp:  "List entries in the stash",
		ShortUsage: "<query>",
		Run:        doStashSearch,
	}
	searchVerb.Add(stashSearch)
}

type stashUploader struct{}

func (stashUploader) prepare(log.Context, *grpc.ClientConn) error { return nil }
func (stashUploader) process(log.Context, string) error           { return nil }

func doStashSearch(ctx log.Context, flags flag.FlagSet) error {
	return grpcutil.Client(ctx, serverAddress, func(ctx log.Context, conn *grpc.ClientConn) error {
		store, err := stashgrpc.Connect(ctx, conn)
		if err != nil {
			return err
		}
		expression := strings.Join(flags.Args(), " ")
		out := ctx.Raw("").Writer()
		expr, err := script.Parse(ctx, expression)
		if err != nil {
			return cause.Explain(ctx, err, "Malformed search query")
		}
		return store.Search(ctx, expr.Query(), func(ctx log.Context, entry *stash.Entity) error {
			proto.MarshalText(out, entry)
			return nil
		})
	}, grpc.WithInsecure())
}
