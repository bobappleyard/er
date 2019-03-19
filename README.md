# ER Modelling for Go

Entity-Relationship (ER) models are a powerful tool for designing software. They are based upon
descriptive diagrams of real-world entities, possessing attributes and connected to one another by
quantity relationships.

This repository maintains a collection of packages that are useful in manipulating a particular
kind of ER model. What makes this approach unusual is that relationships can be annotated with
constraints so that the logical-to-physical transformation is "right first time." The theoretical
basis of this is elucidated on http://www.entitymodelling.org - what you find here is merely an
alternative implementation of the same concepts.

The implementation language of Go is chosen because it balances safety, speed and simplicity, along
with decent code generation facilities. It is fast becoming the top language for writing web-based
infrastructure, making it the perfect basis for creating ER service. The vision behind this project
is reimagining the web-based economy in terms of declaratively-specified distributed data services,
so empowering relatively small teams of programmers to create full-featured business applications
rapidly and with confidence. It is informed by successes in a difference space producing scientific
software.

## Data Modelling: an Exegesis

The conventional approach to ER modelling is as follows. The architect creates a logical model to
describe the application. It is fed into the logical-to-physical transformation, which produces the
attributes required to implement the relationships. Unfortunately this process can create too many
attributes, duplicating information, and so must be manually normalised to fit the intended design.
When the requirements change and the data model requires alteration, the temptation to ignore the
logical model and just change the physical model is hard to resist. Before long, the logical model
bears little to no relationship to the application it was created to describe.

As a result, ER models, and particularly logical models, occupy much too small a niche in the world
of software. However, logical diagrams are a startlingly clear format for presenting ideas. As
communication tool, they enable individuals from a variety of disciplines to agree on a design. As
a design tool, they effect the cenceptualisation and criticism of a data model like no other. They
are software's arch, enabling solid designs that would have otherwise been unthinkable.

Therefore, to restore logical models to their rightful place at the heart of the software
production process, we must remove the manual normalisation step. We achieve this here using an
idea from category theory, the commutativity diagram.

## Entities, Relationships and Attributes

We use *entities* to represent things that exist in the real world that the software is designed
around. A *order* would be an entity, as would a *manager*.

Entities possess *attributes* that describe them.

Entities are connected to other entities by *relationships*. A relationship has a source entity and
a target entity, and we can think of it as pointing from the source to the target.

## Implementing Relationships

Relationships between entities may be subject to *constraints*. Here, a constraint is an assertion
that by following one path (the *diagonal*) from the source of the relationship, we will reach the
same entity as if we followed a different path (the *riser*) from the target of the relationship.
These constraints serve two purposes:

1. They provide a hitherto overlooked tool to aid the designer in describing the constraints that
   real systems are subject to.
2. They provide information to the logical-to-physical transformation, enabling it to create fully
   normalised data models automatically.

Any number of constraints may be present for a given relationship, although in the vast majority of
cases there will only be zero or one. Some archetypical examples are given below.

### Zero Constraints

The zero constraint case is the only situation that most logical-to-physical transformations
support. Here, the identifying attributes of the target entity are effectively copied into the
source entity to implement that relationship.

### One Constraint

When a relationship is subject to a single constraint, only those attributes required to implement
the inverse of the riser are required to be on the source entity. The rest can be inferred from the
implementation of the diagonal.

### Two Constraints

When a relationship is subject to two constraints, the attributes required to implement the
intersection of the inverse of the two risers are required.

This refinement continues. The current implementation uses logical unification to achieve this. As a
result, the source code is simple but extremely subtle.

### Checking Constraints

In order to validate a model, for every relationship, each of the constraints must be checked such
that the asserted equality between following the riser and following the diagonal holds, and any
violations reported.

