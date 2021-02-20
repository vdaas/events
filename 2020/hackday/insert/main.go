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
	"bufio"
	"bytes"
	"context"
	"flag"
	"io"
	"os"
	"strconv"

	"github.com/kpango/glg"
	"google.golang.org/grpc"

	agent "github.com/vdaas/vald-client-go/v1/agent/core"
	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
)

var (
	datasetPath    string
	grpcServerAddr string
)

func init() {
	flag.StringVar(&datasetPath, "path", "jawiki.entity_vectors.100d.txt", "dataset path")
	flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "gRPC server address")
	flag.Parse()
}

func main() {
	dataset, err := load(datasetPath)
	if err != nil {
		glg.Fatal(err)
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		glg.Fatal(err)
	}

	client := vald.NewValdClient(conn)

	for id, vec := range dataset {
		_, err := client.Insert(ctx, &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: vec,
			},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
			},
		})
		if err != nil {
			glg.Fatal(err)
		}
	}

	glg.Info("Start Indexing dataset.")
	_, err = agent.NewAgentClient(conn).CreateIndex(ctx, &payload.Control_CreateIndexRequest{
		PoolSize: uint32(len(dataset) / 100),
	})
	if err != nil {
		glg.Fatal(err)
	}

	glg.Info("Finish Indexing dataset. \n\n")
}

func load(path string) (dataset map[string][]float32, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	head, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	h := bytes.SplitN(head, []byte(" "), 2)
	total, err := strconv.ParseUint(string(h[0]), 10, 64)
	if err != nil {
		return nil, err
	}
	dimension, err := strconv.ParseUint(string(h[1]), 10, 64)
	if err != nil {
		return nil, err
	}
	dataset = make(map[string][]float32, total)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(line) == 0 {
			continue
		}
		row := bytes.Split(line, []byte(" "))
		vec := make([]float32, 0, dimension)
		for _, v := range row[1:] {
			val, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				glg.Warnf("parse failed key:%s num: %s, err %v", string(row[0]), string(v), err)
				continue
			}
			vec = append(vec, float32(val))
		}
		if len(vec) == int(dimension) {
			dataset[string(row[0])] = vec
		}
	}

	return
}
