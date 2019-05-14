package diagram

import (
	"github.com/ajstarks/svgo"
	"github.com/bobappleyard/er/l2p"
	"os"
	"testing"

	. "github.com/bobappleyard/er"
)

func createEntity(m *EntityModel, name string, dependsOn *EntityType) *EntityType {
	res := &EntityType{Name: name}
	m.Types = append(m.Types, res)
	if dependsOn != nil {
		res.DependsOn = addRelationship(res, "parent", dependsOn)
	}
	return res
}

func addAttribute(t *EntityType, name string) *Attribute {
	res := &Attribute{
		Owner: t,
		Name:  name,
		Type:  StringType,
	}
	t.Attributes = append(t.Attributes, res)
	return res
}

func addRelationship(t *EntityType, name string, targ *EntityType) *Relationship {
	res := &Relationship{
		Source: t,
		Target: targ,
		Name:   name,
	}
	t.Relationships = append(t.Relationships, res)
	return res
}

func addConstraint(r *Relationship, diagonal, riser []string) {
	followPath := func(t *EntityType, path []string) []Component {
		res := make([]Component, len(path))
		for i, step := range path {
			for _, r := range t.Relationships {
				if r.Name == step {
					res[i].Rel = r
					t = r.Target
					break
				}
			}
			if res[i].Rel == nil {
				panic("unknown: " + step)
			}
		}
		return res
	}
	r.Constraints = append(r.Constraints, Constraint{
		Diagonal: Diagonal{Components: followPath(r.Source, diagonal)},
		Riser:    Riser{Components: followPath(r.Target, riser)},
	})
}

func TestGenerate(t *testing.T) {
	m := &EntityModel{
		Name: "teams",
	}

	skill := createEntity(m, "skill", nil)
	addAttribute(skill, "name").Identifying = true

	team := createEntity(m, "team", nil)
	addAttribute(team, "name").Identifying = true

	member := createEntity(m, "member", team)
	addAttribute(member, "first_name").Identifying = true
	addAttribute(member, "last_name").Identifying = true

	team_skill := createEntity(m, "team_skill", team)
	team_skill.DependsOn.Identifying = true
	addRelationship(team_skill, "skill", skill).Identifying = true

	learned_skill := createEntity(m, "learned_skill", member)
	learned_skill.DependsOn.Identifying = true
	learned := addRelationship(learned_skill, "learned", team_skill)
	learned.Identifying = true
	addConstraint(learned,
		[]string{"parent", "parent"},
		[]string{"parent"},
	)

	f, err := os.Create("test.svg")
	if err != nil {
		t.Error(err)
		return
	}
	l2p.LogicalToPhysical(m)
	t.Fail()
	Draw(svg.New(f), m)
	f.Close()
}
