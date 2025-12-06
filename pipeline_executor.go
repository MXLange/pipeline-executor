package pipelineexecutor

import "fmt"

// PipelineExecutor manages and executes pipelines consisting of multiple steps.
type PipelineExecutor struct {
	Pipelines map[string]Pipeline
}

// Pipeline represents a sequence of steps to be executed.
type Pipeline struct {
	Name        string
	Key         string
	Steps       []Step
	indexSearch map[string]int
}

// Step represents an individual step within a pipeline.
type Step struct {
	Name          string
	DependsOnStep *int
	Action        func(args ...any) *PipelineError
}

// PipelineError represents an error that occurs during pipeline execution.
type PipelineError struct {
	Err  error
	Step string
}

func (pe *PipelineError) Error() string {
	return pe.Err.Error()
}

// NewPipelineExecutor creates and returns a new PipelineExecutor instance.
func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{
		Pipelines: make(map[string]Pipeline),
	}
}

// AddPipeline adds a new pipeline to the executor.
// If a pipeline with the same key already exists, it will be overwritten.
func (pe *PipelineExecutor) AddPipeline(pipeline Pipeline) {

	pipeline.indexSearch = make(map[string]int)
	for i, step := range pipeline.Steps {
		pipeline.indexSearch[step.Name] = i
	}

	pe.Pipelines[pipeline.Key] = pipeline
}

func (pe *Pipeline) executeStep(step Step, executedSteps map[int]bool, args ...any) *PipelineError {

	if step.DependsOnStep != nil {
		depIndex := *step.DependsOnStep
		if depIndex < 0 || depIndex >= len(pe.Steps) {
			return &PipelineError{
				Err:  fmt.Errorf("invalid dependency index %d for step %s", depIndex, step.Name),
				Step: step.Name,
			}
		}

		if !executedSteps[depIndex] {
			if err := pe.executeStep(pe.Steps[depIndex], executedSteps, args...); err != nil {
				return err
			}
		}
	}

	if err := step.Action(args...); err != nil {
		return err
	}
	return nil
}

func (pe *PipelineExecutor) ExecutePipeline(key, stepToRun string, args ...any) *PipelineError {
	pipeline, exists := pe.Pipelines[key]
	if !exists {
		return &PipelineError{
			Err:  fmt.Errorf("pipeline with key %s not found", key),
			Step: "",
		}
	}

	startIndex := 0

	if stepToRun != "" {
		i, found := pipeline.indexSearch[stepToRun]
		if !found {
			return &PipelineError{
				Err:  fmt.Errorf("step with name %s not found in pipeline %s", stepToRun, key),
				Step: "",
			}
		}

		startIndex = i
	}

	executedSteps := make(map[int]bool)

	for i, step := range pipeline.Steps {

		if i < startIndex {
			continue
		}

		if err := pipeline.executeStep(step, executedSteps, args...); err != nil {
			return err
		}
		executedSteps[i] = true
	}

	return nil
}
