Path Language
=============

An important part of this package is a DSL for describing relationships. This
DSL is referred to as the *path language*. Tools exist for parsing paths and
for evaluating paths againsta a model. What follows is a description of the 
syntax and an exploration of some of the surrounding concepts, along with links
to the tools that are based upon paths.

Background
==========

Before getting to the syntax, it is important to understand what it is we are
trying to achieve. A relationship is a kind of mapping from one set of entities
to another set of entities. We might therefore say that there is a *source set*
and a *target set*.

Paths generalise this a little, because they can also describe mappings that
are not between entity types at all. The other types of set that may be source
or target sets of paths are values and the absolute.

Syntax
======

The syntax of paths is intended to concisely describe the ways in which
relationships are usually formed in an ER model. Due to the concision, there
are potential ambiguities, but these are resolvable.

Paths have the following syntax:

    Path         = Union
    Union        = Intersection | Intersection '|' Union
    Intersection = Join | Join '&' Intersection
    Join         = Inverse | Inverse '/' Join
    Inverse      = Term | '~' Term
    Term         = '*' | IDENT | STR | '(' Path ')'

Some examples:

    a/b/c
    emp_dept_id/~dept_id&emp_no/~no
    */mp&'Jim Hacker'/~name

Terms
-----

A quoted strings refers to the value of that string unquoted. 

Combinations
------------

Paths can be combined into more complicated paths, so as to describe the

Inverse
-------

Paths may be inverted. This has the effect of swapping the source set and the
target set.

This has some interesting implications in the case of paths that refer to
attributes. If a path `a` refers to an attribute, i.e. it has an entity type
as its source set and a value type as its target set, then `~a` refers to a
filter over attribute values. This means that it maps to all the entities that
contain an attribute named `a` that have a particular value. This is very
useful in describing relationships. We can use join to say, for example,
`target_name/~name`, which means "find the value of `target_name` for the current entity