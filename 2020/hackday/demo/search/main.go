//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/kpango/glg"
	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"google.golang.org/grpc"
)

var (
	query          string
	topk           uint64
	grpcServerAddr string
)

func init() {
	flag.StringVar(&query, "q", "", "search by id key")
	flag.Uint64Var(&topk, "k", 10, "top k")
	flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "gRPC server address")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		glg.Fatal(err)
	}

	resp, err := vald.NewValdClient(conn).SearchByID(ctx, &payload.Search_IDRequest{
		Id: query,
		Config: &payload.Search_Config{
			Num:     uint32(topk),
			Radius:  -1,
			Epsilon: 0.1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, result := range resp.GetResults() {
		fmt.Fprintf(w, "id: %s \t distance: %f\n", result.GetId(), result.GetDistance())
	}
	w.Flush()
}
