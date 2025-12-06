package pipelineexecutor

import "fmt"

type PipelineExecutor struct {
	Pipelines map[string]Pipeline
}

type Pipeline struct {
	Name  string
	Key   string
	Steps []Step
}

type Step struct {
	Name          string
	DependsOnStep *int
	Action        func(args ...any) *PipelineError
}

type PipelineError struct {
	Error error
	Step  string
}

func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{
		Pipelines: make(map[string]Pipeline),
	}
}

func (pe *PipelineExecutor) AddPipeline(pipeline Pipeline) {
	pe.Pipelines[pipeline.Key] = pipeline
}

func (pe *Pipeline) executeStep(step Step, executedSteps map[int]bool, args ...any) *PipelineError {

	if step.DependsOnStep != nil {
		depIndex := *step.DependsOnStep
		if depIndex < 0 || depIndex >= len(pe.Steps) {
			return &PipelineError{
				Error: fmt.Errorf("invalid dependency index %d for step %s", depIndex, step.Name),
				Step:  step.Name,
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
			Error: fmt.Errorf("pipeline with key %s not found", key),
			Step:  "",
		}
	}

	totalSteps := len(pipeline.Steps)
	executedSteps := make(map[int]bool)

	for i, step := range pipeline.Steps {
		if stepToRun != "" && step.Name != stepToRun {
			totalSteps--
			continue
		}

		if err := pipeline.executeStep(step, executedSteps, args...); err != nil {
			return err
		}
		executedSteps[i] = true
	}

	if totalSteps == 0 {
		return &PipelineError{
			Error: fmt.Errorf("step with name %s not found", stepToRun),
			Step:  "",
		}
	}

	return nil
}
