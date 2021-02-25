package app

import (
	"context"

	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ConfigurePipeline(ctx context.Context, p *Pipeline) error {
	if err := p.Deps.WorkflowRun.Own(ctx, p); err != nil {
		return err
	}

	if err := ConfigureConditions(ctx, p.Conditions); err != nil {
		return err
	}

	if err := ConfigureTasks(ctx, p.Tasks); err != nil {
		return err
	}

	p.Object.Spec.Tasks = make([]tektonv1beta1.PipelineTask, 0, len(p.Tasks.List))

	for i, t := range p.Tasks.List {
		ws := p.Deps.WorkflowRun.Object.Spec.Workflow.Steps[i]
		ms := ModelStep(p.Deps.WorkflowRun, ws)

		pt := tektonv1beta1.PipelineTask{
			Name: ms.Hash().HexEncoding(),
			TaskRef: &tektonv1beta1.TaskRef{
				Name: t.Key.Name,
			},
			RunAfter: make([]string, len(ws.DependsOn)),
		}

		for i, dep := range ws.DependsOn {
			pt.RunAfter[i] = ModelStepFromName(p.Deps.WorkflowRun, dep).Hash().HexEncoding()
		}

		if cond, ok := p.Conditions.GetByStepName(ws.Name); ok {
			pt.Conditions = []tektonv1beta1.PipelineTaskCondition{
				{ConditionRef: cond.Key.Name},
			}
		}

		p.Object.Spec.Tasks = append(p.Object.Spec.Tasks, pt)
	}

	return nil
}

func ApplyPipeline(ctx context.Context, cl client.Client, deps *WorkflowRunDeps) (*Pipeline, error) {
	p := NewPipeline(deps)

	if _, err := p.Load(ctx, cl); err != nil {
		return nil, err
	}

	p.LabelAnnotateFrom(ctx, deps.WorkflowRun.Object.ObjectMeta)

	if err := ConfigurePipeline(ctx, p); err != nil {
		return nil, err
	}

	if err := p.Persist(ctx, cl); err != nil {
		return nil, err
	}

	return p, nil
}
