package l2p

import (
	"github.com/bobappleyard/er"
)

func missingAttributes(rels map[*er.Relationship]resolvedPath, r *er.Relationship) []*er.Attribute {
	p := rels[r]
	source, target := p.route()
	// Start out with all attributes defined on the source
	var attrs resolvedPath = resolvedInverse{resolvedAttribute{source.Attributes[0]}}
	for _, a := range source.Attributes[1:] {
		attrs = resolvedIntersection{attrs, resolvedInverse{resolvedAttribute{a}}}
	}
	ctx := &analysisCtx{rels: rels}
	ctx = &analysisCtx{ctx, rels, resolvedJoin{
		resolvedValue{""},
		resolvedJoin{attrs, p},
	}}
	// Attempt to find the attributes in the key of the target using the path provided.
	k := key(target)
	res := make([]*er.Attribute, 0, len(k))
	for _, a := range k {
		if ctx.provides(a) {
			continue
		}
		// The path does not imply this attribute, so it needs to be included.
		res = append(res, a)
	}
	return res
}

type analysisCtx struct {
	prev *analysisCtx
	rels map[*er.Relationship]resolvedPath
	path resolvedPath
}

func (c *analysisCtx) provides(a *er.Attribute) bool {
	// Either the attribute is explicitly declared as existing...
	if c.path.provides(c.prev, false, a) {
		return true
	}
	if a == attributeValue {
		return false
	}
	// ...or the key is present, which implies all attributes
	_, target := c.path.route()
	for _, b := range key(target) {
		if a == b {
			// We already tried to find the attribute and failed, above.
			return false
		}
		if !c.path.provides(c.prev, false, b) {
			return false
		}
	}
	return true
}

func (p absolute) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return false
}

func (p resolvedEntityType) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return false
}

func (p resolvedRelationship) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return ctx.rels[p.r].provides(ctx, inv, a)
}

func (p resolvedAttribute) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	left, right := attributeValue, p.a
	if inv {
		left, right = right, left
	}
	return a == left && ctx.provides(right)
}

func (p resolvedValue) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return !inv && a == attributeValue
}

func (p resolvedInverse) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return p.p.provides(ctx, !inv, a)
}

func (p resolvedJoin) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	left, right := p.left, p.right
	if inv {
		left, right = resolvedInverse{right}, left
	}
	return right.provides(&analysisCtx{ctx, ctx.rels, left}, inv, a)
}

func (p resolvedIntersection) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return p.left.provides(ctx, inv, a) || p.right.provides(ctx, inv, a)
}

func (p resolvedUnion) provides(ctx *analysisCtx, inv bool, a *er.Attribute) bool {
	return p.left.provides(ctx, inv, a) && p.right.provides(ctx, inv, a)
}
