The Model
=========

Entity-Relationship models represent data using entity types that 
are connected by relationships. These relationships are typically
implemented by attributes subject to constraints. In this package,
the constraints are expressed in terms of paths.

This document specifies ER models so that it is possible to create your own after reading it.

Logical and Physical Models
===========================

ER models may be treated in different ways 

Entities and Attributes
=======================

Entity Type
-----------

Attribute
---------

Identity
--------

Relationships and Constraints 
=============================

A relationship connects entities. This places constraints over what
entities are legal based on other entities in the model. These constraints, when properly asserted, prevent bad data from entering your system and breaking your applications.

Every relationship is defined on a *source entity* and points to a 
*target entity*.

    entity_type {
        name: "a"

        relationship {
            name: "f"
            type: "b"
        }
    }

    entity_type {
        name: "b"
    }

In this example, the relationship `f` has source entities of type `a` and target entities of type `b`.

First-Order Relationship
------------------------

A first-order relationship is one that is implemented entirely in
terms of attributes on the source entity.

Commutativity Constraint
------------------------

Pullback Constraint
-------------------

Recursive Models
================

