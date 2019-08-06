package l2p

import (
	"github.com/bobappleyard/er"
	"github.com/bobappleyard/er/util/top"
)

// Sort relationships so that, for each relationship, the relationships upon
// which it depends have already been implemented.
//
// * Relationships depend upon their target entity type.
//
// * Relationships depend upon any relationships that are mentioned in their path.
//
// * Entity types depend upon all identifying relationships.
//
// This dependency applies transitively, so we use a topological sort to
// accomplish this.
func sortRels(rels map[*er.Relationship]resolvedPath) ([]*er.Relationship, error) {
	var g top.Graph
	for r, rp := range rels {
		g.Link(r.Target, r)
		for _, r := range r.Source.Relationships {
			if !r.Identifying {
				continue
			}
			g.Link(r, r.Source)
		}
		for _, s := range rp.appendRels(nil) {
			g.Link(s, r)
		}
	}
	sorted, err := g.Sort()
	if err != nil {
		return nil, err
	}
	var res []*er.Relationship
	for _, r := range sorted {
		if r, ok := r.(*er.Relationship); ok {
			res = append(res, r)
		}
	}
	return res, nil
}

func (p absolute) appendRels(rels []*er.Relationship) []*er.Relationship {
	return rels
}

func (p resolvedValue) appendRels(rels []*er.Relationship) []*er.Relationship {
	return rels
}

func (p resolvedEntityType) appendRels(rels []*er.Relationship) []*er.Relationship {
	return rels
}

func (p resolvedRelationship) appendRels(rels []*er.Relationship) []*er.Relationship {
	return append(rels, p.r)
}

func (p resolvedAttribute) appendRels(rels []*er.Relationship) []*er.Relationship {
	return rels
}

func (p resolvedInverse) appendRels(rels []*er.Relationship) []*er.Relationship {
	return p.p.appendRels(rels)
}

func (p resolvedJoin) appendRels(rels []*er.Relationship) []*er.Relationship {
	return p.left.appendRels(p.right.appendRels(rels))
}

func (p resolvedIntersection) appendRels(rels []*er.Relationship) []*er.Relationship {
	return p.left.appendRels(p.right.appendRels(rels))
}

func (p resolvedUnion) appendRels(rels []*er.Relationship) []*er.Relationship {
	return p.left.appendRels(p.right.appendRels(rels))
}
