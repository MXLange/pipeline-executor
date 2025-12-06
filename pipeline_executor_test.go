package pipelineexecutor

import (
	"fmt"
	"testing"
)

func TestPipelineExecutor(t *testing.T) {
	pe := NewPipelineExecutor()

	arg1 := "initial"
	arg2 := 10

	step1 := Step{
		Name: "Step 1",
		Action: func(args ...any) *PipelineError {
			a1 := args[0].(*string)
			a2 := args[1].(*int)

			*a1 = "modified"
			*a2 = *a2 * 2

			return nil
		},
	}

	dependableStep := 0

	step2 := Step{
		Name: "Step 2",
		Action: func(args ...any) *PipelineError {

			a1 := args[0].(*string)
			a2 := args[1].(*int)

			*a1 = *a1 + " again"
			*a2 = *a2 * 2
			return nil
		},
		DependsOnStep: &dependableStep,
	}

	pipeline := Pipeline{
		Name:  "Test Pipeline",
		Key:   "test_pipeline",
		Steps: []Step{step1, step2},
	}

	pe.AddPipeline(pipeline)

	err := pe.ExecutePipeline("test_pipeline", "", &arg1, &arg2)
	if err != nil {
		panic(err.Error())
	}

	expectedArg1 := "modified again"
	expectedArg2 := 40

	if arg1 != expectedArg1 {
		t.Errorf("Expected arg1 to be %s, but got %s", expectedArg1, arg1)
	}

	if arg2 != expectedArg2 {
		t.Errorf("Expected arg2 to be %d, but got %d", expectedArg2, arg2)
	}

	arg1 = "initial"
	arg2 = 10

	err = pe.ExecutePipeline("test_pipeline", "Step 2", &arg1, &arg2)
	if err != nil {
		panic(err.Error())
	}

	expectedArg1 = "modified again"
	expectedArg2 = 40

	if arg1 != expectedArg1 {
		t.Errorf("Expected arg1 to be %s, but got %s", expectedArg1, arg1)
	}

	if arg2 != expectedArg2 {
		t.Errorf("Expected arg2 to be %d, but got %d", expectedArg2, arg2)
	}

	err = pe.ExecutePipeline("test_pipeline", "Step Not Found", &arg1, &arg2)
	if err == nil {
		t.Errorf("Expected error for non-existent step, but got nil")
	}

	if err.Error() != "step with name Step Not Found not found in pipeline test_pipeline" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	step2 = Step{
		Name: "Step 2",
		Action: func(args ...any) *PipelineError {

			a1 := args[0].(*string)
			a2 := args[1].(*int)

			*a1 = *a1 + " again"
			*a2 = *a2 * 2
			return nil
		},
	}

	pipeline = Pipeline{
		Name:  "Test Pipeline",
		Key:   "test_pipeline",
		Steps: []Step{step1, step2},
	}

	pe.AddPipeline(pipeline)

	arg1 = "initial"
	arg2 = 10

	err = pe.ExecutePipeline("test_pipeline", "Step 2", &arg1, &arg2)
	if err != nil {
		panic(err.Error())
	}

	expectedArg1 = "initial again"
	expectedArg2 = 20

	if arg1 != expectedArg1 {
		t.Errorf("Expected arg1 to be %s, but got %s", expectedArg1, arg1)
	}

	if arg2 != expectedArg2 {
		t.Errorf("Expected arg2 to be %d, but got %d", expectedArg2, arg2)
	}

	err = pe.ExecutePipeline("test_pipeline", "Step Not Found", &arg1, &arg2)
	if err == nil {
		t.Errorf("Expected error for non-existent step, but got nil")
	}

	if err.Error() != "step with name Step Not Found not found in pipeline test_pipeline" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

}

func TestPipelineExecutorErrorRetry(t *testing.T) {
	pe := NewPipelineExecutor()

	arg := 10
	var errFlag *bool = new(bool)
	*errFlag = true

	step1 := Step{
		Name: "Step 1",
		Action: func(args ...any) *PipelineError {
			a2 := args[0].(*int)

			*a2 = *a2 * 2

			return nil
		},
	}

	step2 := Step{
		Name: "Step 2",
		Action: func(args ...any) *PipelineError {

			a2 := args[0].(*int)

			if *errFlag {
				*errFlag = false
				return &PipelineError{
					Err: fmt.Errorf("an error occurred in Step 2"),
				}
			}

			*a2 = *a2 * 2
			return nil
		},
	}

	step3 := Step{
		Name: "Step 3",
		Action: func(args ...any) *PipelineError {

			a2 := args[0].(*int)

			*a2 = *a2 * 2
			return nil
		},
	}

	pipeline := Pipeline{
		Name:  "Test Pipeline",
		Key:   "test_pipeline",
		Steps: []Step{step1, step2, step3},
	}

	pe.AddPipeline(pipeline)

	err := pe.ExecutePipeline("test_pipeline", "", &arg)
	if err != nil {
		fmt.Printf("First execution error: %s\n", err.Error())

		expectedArg := 20

		if arg != expectedArg {
			t.Errorf("Expected arg to be %d, but got %d", expectedArg, arg)
		}

		arg = 10

		err = pe.ExecutePipeline("test_pipeline", err.Step, &arg)
		if err != nil {
			t.Errorf("Retry failed with error: %s", err.Error())
		}
	}

	expectedArg := 40

	if arg != expectedArg {
		t.Errorf("Expected arg to be %d, but got %d", expectedArg, arg)
	}

}
