package l2p

import (
	"errors"
	"fmt"
	"github.com/bobappleyard/er"
)

// Errors that can be returned by LogicalToPhysical
var (
	ErrNoPath        = errors.New("no path through model")
	ErrAmbiguousPath = errors.New("ambiguous path through model")
)

// LogicalToPhysical ensures that the relationships in a model are fully
// implemented.
func LogicalToPhysical(m *er.EntityModel) error {
	// Step 1: resolve names in paths
	resolved := map[*er.Relationship]resolvedPath{}
	for _, e := range m.Types {
		for _, r := range e.Relationships {
			rp, err := resolveRelationship(m, r)
			if err != nil {
				return err
			}
			resolved[r] = rp
		}
	}
	// Step 2: sort the relationships topologically
	rels, err := sortRels(resolved)
	if err != nil {
		return err
	}
	// Step 3: insert attributes that are missing from the scope paths
	for _, r := range rels {
		for _, a := range missingAttributes(resolved, r) {
			nm := r.Name + "_" + a.Name
			r.Source.Attributes = append(r.Source.Attributes, &er.Attribute{
				Name:        nm,
				Owner:       r.Source,
				Type:        a.Type,
				Identifying: r.Identifying,
			})
			if r.Path != "" {
				r.Path += "&"
			}
			r.Path += fmt.Sprintf("%s/~%s", nm, a.Name)
		}
		rp, _ := resolveRelationship(m, r)
		resolved[r] = rp
	}
	return nil
}

var (
	absoluteType   = &er.EntityType{Name: "*"}
	valueType      = &er.EntityType{Name: "$"}
	attributeValue = &er.Attribute{Owner: valueType}
)

type resolvedPath interface {
	route() (source, target *er.EntityType)
	appendRels(rels []*er.Relationship) []*er.Relationship
	provides(ctx *analysisCtx, a *er.Attribute) bool
	inverseProvides(ctx *analysisCtx, a *er.Attribute) bool
}

type absolute struct{}

type resolvedValue struct {
	v string
}

type resolvedEntityType struct {
	e *er.EntityType
}

type resolvedRelationship struct {
	r *er.Relationship
}

type resolvedAttribute struct {
	a *er.Attribute
}

type resolvedInverse struct {
	p resolvedPath
}

type resolvedJoin struct {
	left, right resolvedPath
}

type resolvedIntersection struct {
	left, right resolvedPath
}

type resolvedUnion struct {
	left, right resolvedPath
}

func (p absolute) String() string {
	return "*"
}

func (p resolvedValue) String() string {
	return "'" + p.v + "'"
}

func (p resolvedAttribute) String() string {
	return p.a.Name
}

func (p resolvedEntityType) String() string {
	return p.e.Name
}

func (p resolvedRelationship) String() string {
	return p.r.Name
}

func (p resolvedInverse) String() string {
	return fmt.Sprintf("~%s", p.p)
}

func (p resolvedJoin) String() string {
	return fmt.Sprintf("(%s)/(%s)", p.left, p.right)
}

func (p resolvedIntersection) String() string {
	return fmt.Sprintf("(%s)&(%s)", p.left, p.right)
}

func (p resolvedUnion) String() string {
	return fmt.Sprintf("(%s)|(%s)", p.left, p.right)
}
