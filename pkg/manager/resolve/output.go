package resolve

import (
	"context"

	exprmodel "github.com/puppetlabs/relay-core/pkg/expr/model"
	"github.com/puppetlabs/relay-core/pkg/expr/resolve"
	"github.com/puppetlabs/relay-core/pkg/model"
)

type OutputTypeResolver struct {
	m model.StepOutputGetterManager
}

var _ resolve.OutputTypeResolver = &OutputTypeResolver{}

func (otr *OutputTypeResolver) ResolveOutput(ctx context.Context, from, name string) (interface{}, error) {
	so, err := otr.m.Get(ctx, from, name)
	if err == model.ErrNotFound {
		return nil, &exprmodel.OutputNotFoundError{From: from, Name: name}
	} else if err != nil {
		return nil, err
	}

	return so.Value, nil
}

func NewOutputTypeResolver(m model.StepOutputGetterManager) *OutputTypeResolver {
	return &OutputTypeResolver{
		m: m,
	}
}
