# 📦 pipeline-executor

A small Go library for creating and executing **pipelines composed of
multiple steps**, supporting:

-   **Sequential execution**
-   **Step dependencies**
-   **Passing arguments by reference**
-   **Automatic retry when a step fails**
-   **Partial execution** (start from any step)

Define steps as functions using `args ...any`, register a pipeline, and
the executor manages order, dependencies, and errors.

---

## 🚀 Installation

```bash
go get github.com/MXLange/pipeline-executor
```

---

## 🧱 Core Concepts

### **PipelineExecutor**

Manages pipeline registration and execution.

### **Pipeline**

Contains: - `Name` - `Key` - `Steps` - Internal index mapping to locate
steps efficiently

### **Step**

A pipeline step has: - `Name` - `DependsOnStep` (optional) -
`Action(args ...any)` → `*PipelineError`

### **PipelineError**

Structure containing: - `Err` (original Go error) - `Step` (step name
where it failed)

---

## 📌 Example Usage

```go
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
            fmt.Println("Step1:", a.Value, b.Number)
            a.Value = "modified"
            return nil
        },
    }

    step2 := pipelineexecutor.Step{
        Name: "Step 2",
        Action: func(args ...any) *pipelineexecutor.PipelineError {
            a := args[0].(*ExStruct)
            b := args[1].(*IntExStruct)
            fmt.Println("Step2:", a.Value, b.Number)
            b.Number *= 2
            a.Value += " again"
            return nil
        },
    }

    pipeline := pipelineexecutor.Pipeline{
        Name:  "Example Pipeline",
        Key:   "example_pipeline",
        Steps: []pipelineexecutor.Step{step1, step2},
    }

    pe.AddPipeline(pipeline)

    err := pe.ExecutePipeline("example_pipeline", "", arg, intArg)
    if err != nil {
        panic(err)
    }

    fmt.Println("Final:", arg.Value, intArg.Number)
}
```

---

## 🔁 Automatic Retry on Failure

```go
err := pe.ExecutePipeline("test_pipeline", "", &arg)
if err != nil {
    pe.ExecutePipeline("test_pipeline", err.Step, &arg)
}
```

---

## 🎯 Partial Execution

```go
pe.ExecutePipeline("example_pipeline", "Step 2", &arg1, &arg2)
```

---

## 🔗 Step Dependencies

```go
dependIndex := 0
step2 := Step{
    Name: "Step 2",
    DependsOnStep: &dependIndex,
    Action: func(args ...any) *PipelineError {
        return nil
    },
}
```

---

## 🧪 Testing

Run tests:

```bash
go test ./...
```

---

## 📄 License

**MIT**
