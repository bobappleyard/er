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
	res := make([]*er.Attribute, 0, len(target.Attributes))
	for _, a := range target.Attributes {
		if !a.Identifying || ctx.provides(a) {
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
	if c.path.provides(c.prev, a) {
		return true
	}
	if a == attributeValue {
		return false
	}
	// ...or the key is present, which implies all attributes
	for _, b := range c.key() {
		if a == b {
			// We already tried to find the attribute and failed, above.
			return false
		}
		if !c.path.provides(c.prev, b) {
			return false
		}
	}
	return true
}

func (c *analysisCtx) key() []*er.Attribute {
	_, e := c.path.route()
	key := make([]*er.Attribute, 0, len(e.Attributes))
	for _, a := range e.Attributes {
		if !a.Identifying {
			continue
		}
		key = append(key, a)
	}
	return key
}

func (p absolute) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return false
}

func (p absolute) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return false
}

func (p resolvedEntityType) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return false
}

func (p resolvedEntityType) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return false
}

func (p resolvedRelationship) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return ctx.rels[p.r].provides(ctx, a)
}

func (p resolvedRelationship) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return ctx.rels[p.r].inverseProvides(ctx, a)
}

func (p resolvedAttribute) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return a == attributeValue && ctx.provides(p.a)
}

func (p resolvedAttribute) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return a == p.a && ctx.provides(attributeValue)
}

func (p resolvedValue) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return a == attributeValue
}

func (p resolvedValue) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return false
}

func (p resolvedInverse) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.p.inverseProvides(ctx, a)
}

func (p resolvedInverse) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.p.provides(ctx, a)
}

func (p resolvedJoin) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.right.provides(&analysisCtx{ctx, ctx.rels, p.left}, a)
}

func (p resolvedJoin) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.left.inverseProvides(&analysisCtx{ctx, ctx.rels, resolvedInverse{p.right}}, a)
}

func (p resolvedIntersection) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.left.provides(ctx, a) || p.right.provides(ctx, a)
}

func (p resolvedIntersection) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.left.inverseProvides(ctx, a) || p.right.inverseProvides(ctx, a)
}

func (p resolvedUnion) provides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.left.provides(ctx, a) && p.right.provides(ctx, a)
}

func (p resolvedUnion) inverseProvides(ctx *analysisCtx, a *er.Attribute) bool {
	return p.left.inverseProvides(ctx, a) && p.right.inverseProvides(ctx, a)
}
