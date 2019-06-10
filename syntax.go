package er

import (
	"errors"
	"fmt"
)

// EntityModel represents a collection of entity types and how they relate.
type EntityModel struct {
	Name  string
	Types []*EntityType
}

// EntityType represents an entity type.
type EntityType struct {
	Name          string          `rsf:"name"`
	Attributes    []*Attribute    `rsf:"attribute"`
	Relationships []*Relationship `rsf:"relationship"`
	DependsOn     *Relationship
}

// Attribute represents an attribute.
type Attribute struct {
	Name        string        `rsf:"name"`
	Type        AttributeType `rsf:"type"`
	Identifying bool          `rsf:"identifying"`
	Owner       *EntityType
}

// AttributeType represents an attribute type.
type AttributeType byte

// Supported attribute types.
const (
	InvalidType AttributeType = iota
	StringType
	IntType
	FloatType
)

// Relationship represents a relationship.
type Relationship struct {
	Name           string           `rsf:"name"`
	TargetName     string           `rsf:"type_name"`
	Constraints    []Constraint     `rsf:"constraint"`
	Identifying    bool             `rsf:"identifying"`
	Implementation []Implementation `rsf:"implementation"`
	Source, Target *EntityType
}

// Constraint represnts a constraint over a relationship.
type Constraint struct {
	Diagonal Diagonal `rsf:"diagonal"`
	Riser    Riser    `rsf:"riser"`
}

type Diagonal struct {
	Components []Component `rsf:"component"`
}

type Riser struct {
	Components []Component `rsf:"component"`
}

// Component represents part of a path.
type Component struct {
	RelName string `rsf:"rel_name"`
	Rel     *Relationship
}

// Implementation represents part of an attribute filter.
type Implementation struct {
	Source, Target *Attribute
	BasePath       []Component
}

func (a *Attribute) String() string {
	return fmt.Sprintf("%s.%s", a.Owner.Name, a.Name)
}

func (t *EntityType) String() string {
	return t.Name
}

func (a *Relationship) String() string {
	return fmt.Sprintf("%s.%s", a.Source.Name, a.Name)
}

// ErrInvalidAttribute ...
var (
	ErrInvalidAttribute = errors.New("invalid attribute")
	ErrInvalidRecord    = errors.New("invalid record type")
	ErrMissingAttribute = errors.New("missing attribute")
	ErrDuplicateKey     = errors.New("duplicate key")
	ErrMissingEntity    = errors.New("entity not found by key")
	ErrImmutableSet     = errors.New("attempting to modify immutable set")
)
