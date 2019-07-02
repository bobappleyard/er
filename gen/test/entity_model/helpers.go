package entity_model

func (e EntityType) Relationships() setOfRelationship {
	t := e.model.Relationship
	return t.Where(t.SourceName.Eq(e.Name))
}

func (e EntityType) IncomingRelationships() setOfRelationship {
	t := e.model.Relationship
	return t.Where(t.TargetName.Eq(e.Name))
}

func (e EntityType) Attributes() setOfAttribute {
	t := e.model.Attribute
	return t.Where(t.OwnerName.Eq(e.Name))
}

func (e EntityType) Dependency() setOfDependency {
	m := e.model
	return m.Dependency.Where(m.Dependency.EntityTypeName.Eq(e.Name))
}

func (e Relationship) Dependency() setOfDependency {
	m := e.model
	q := m.Dependency.EntityTypeName.Eq(e.SourceName)
	q = q.And(m.Dependency.RelationshipName.Eq(e.Name))
	return m.Dependency.Where(q)
}

func (m *Model) RootTypes() setOfEntityType {
	var res setOfEntityType
	m.EntityType.ForEach(func(e EntityType) error {
		if e.Dependency().Count() == 0 {
			res = res.Union(m.EntityType.Where(e.identity()))
		}
		return nil
	})
	return res
}

func (e EntityType) Dependents() setOfEntityType {
	var res setOfEntityType
	e.IncomingRelationships().ForEach(func(r Relationship) error {
		r.Dependency().ForEach(func(d Dependency) error {
			res = res.Union(d.queryForEntityType())
			return nil
		})
		return nil
	})
	return res
}
