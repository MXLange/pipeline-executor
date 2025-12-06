package main

import (
	"fmt"

	pipelineexecutor "github.com/MXLange/pipeline-executor"
)

type ExStruct struct {
	Value string
}

func main() {

	pe := pipelineexecutor.NewPipelineExecutor()

	arg := &ExStruct{Value: "initial"}

	var errFlag bool = true

	step1 := pipelineexecutor.Step{
		Name: "Step 1",
		Action: func(args ...any) *pipelineexecutor.PipelineError {
			a := args[0].(*ExStruct)
			fmt.Printf("step 1 value is -> %s\n", a.Value)

			a.Value = "modified"
			return nil
		},
	}

	step2 := pipelineexecutor.Step{
		Name: "Step 2",
		Action: func(args ...any) *pipelineexecutor.PipelineError {
			a := args[0].(*ExStruct)

			fmt.Printf("step 2 value is -> %s\n", a.Value)

			a.Value = a.Value + " again"
			return nil
		},
	}

	step3 := pipelineexecutor.Step{
		Name: "Step 3",
		Action: func(args ...any) *pipelineexecutor.PipelineError {
			a := args[0].(*ExStruct)

			fmt.Printf("step 3 value is -> %s\n", a.Value)

			if errFlag {
				errFlag = false
				return &pipelineexecutor.PipelineError{
					Err:  fmt.Errorf("an error occurred in Step 3"),
					Step: "Step 3",
				}
			}

			a.Value = a.Value + " again"
			return nil
		},
	}

	pipeline := pipelineexecutor.Pipeline{
		Name:  "example Pipeline",
		Key:   "example_pipeline",
		Steps: []pipelineexecutor.Step{step1, step2, step3},
	}

	pe.AddPipeline(pipeline)

	err := pe.ExecutePipeline("example_pipeline", "", arg, errFlag)
	if err != nil {
		arg.Value = "initial"
		fmt.Printf("Pipeline execution failed at step: %s with error: %s\n", err.Step, err.Error())
	}

	fmt.Printf("executing pipeline again...\n")

	err = pe.ExecutePipeline("example_pipeline", err.Step, arg, errFlag)
	if err != nil {
		fmt.Printf("Pipeline execution failed at step: %s with error: %s\n", err.Step, err.Error())
	}

	fmt.Printf("final value is -> %s\n", arg.Value)
}
