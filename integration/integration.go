// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"github.com/pingcap/log"
	"github.com/pingcap/ticdc/integration/framework"
	canal2 "github.com/pingcap/ticdc/integration/framework/canal"
	"github.com/pingcap/ticdc/integration/tests/avro"
	"github.com/pingcap/ticdc/integration/tests/canal"
	"go.uber.org/zap"

	avro2 "github.com/pingcap/ticdc/integration/framework/avro"

	"go.uber.org/zap/zapcore"
)

var testProtocol = flag.String("protocol", "avro", "the protocol we want to test: avro or canal")
var dockerComposeFile = flag.String("docker-compose-file", "", "the path of the Docker-compose yml file")

func testAvro() {
	testCases := []framework.Task{
		avro.NewSimpleCase(),
		avro.NewDeleteCase(),
		avro.NewManyTypesCase(),
		avro.NewUnsignedCase(),
		avro.NewCompositePKeyCase(),
		avro.NewAlterCase(), // this case is slow, so put it last
	}

	log.SetLevel(zapcore.DebugLevel)
	env := avro2.NewKafkaDockerEnv(*dockerComposeFile)
	env.Setup()

	for i := range testCases {
		env.RunTest(testCases[i])
		if i < len(testCases)-1 {
			env.Reset()
		}
	}

	env.TearDown()
}

func testCanal() {
	testCases := []framework.Task{
		canal.NewSimpleCase(),
		canal.NewDeleteCase(),
		canal.NewManyTypesCase(),
		//canal.NewUnsignedCase(), //now canal adapter can not deal with unsigned int greater than int max
		canal.NewCompositePKeyCase(),
		//canal.NewAlterCase(), // basic implementation can not grantee ddl dml sequence, so can not pass
	}

	log.SetLevel(zapcore.DebugLevel)
	env := canal2.NewKafkaDockerEnv(*dockerComposeFile)
	env.Setup()

	for i := range testCases {
		env.RunTest(testCases[i])
		if i < len(testCases)-1 {
			env.Reset()
		}
	}

	env.TearDown()
}

func main() {
	flag.Parse()
	if *testProtocol == "avro" {
		testAvro()
	} else if *testProtocol == "canal" {
		testCanal()
	} else
	{
		log.Fatal("Unknown sink protocol", zap.String("protocol", *testProtocol))
	}
}
