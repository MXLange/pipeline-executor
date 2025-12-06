package main

import (
	"fmt"

	pipelineexecutor "github.com/MXLange/pipeline-executor"
)

type ExStruct struct {
	Value string
}

type IntExStruct struct {
	Number int
}

func main() {

	pe := pipelineexecutor.NewPipelineExecutor()

	arg := &ExStruct{Value: "initial"}
	intArg := &IntExStruct{Number: 5}

	step1 := pipelineexecutor.Step{
		Name: "Step 1",
		Action: func(args ...any) *pipelineexecutor.PipelineError {
			a := args[0].(*ExStruct)
			b := args[1].(*IntExStruct)
			fmt.Printf("initial value is -> %s\n", a.Value)
			fmt.Printf("initial number is -> %d\n", b.Number)

			a.Value = "modified"
			return nil
		},
	}

	step2 := pipelineexecutor.Step{
		Name: "Step 2",
		Action: func(args ...any) *pipelineexecutor.PipelineError {
			a := args[0].(*ExStruct)
			b := args[1].(*IntExStruct)

			fmt.Printf("first modified value is -> %s\n", a.Value)
			fmt.Printf("not modified number is -> %d\n", b.Number)

			b.Number = b.Number * 2

			a.Value = a.Value + " again"
			return nil
		},
	}

	pipeline := pipelineexecutor.Pipeline{
		Name:  "example Pipeline",
		Key:   "example_pipeline",
		Steps: []pipelineexecutor.Step{step1, step2},
	}

	pe.AddPipeline(pipeline)

	err := pe.ExecutePipeline("example_pipeline", "", arg, intArg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("final value is -> %s\n", arg.Value)
	fmt.Printf("final number is -> %d\n", intArg.Number)
}
