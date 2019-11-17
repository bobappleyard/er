Introduction
============

Getting Started
---------------

To use this package currently requires a recent version of the Go
programming language. If you do not already have this installed, consult [its website](https://www.golang.org).

Then execute the following command:

    go get github.com/bobappleyard/er

The package is now on your system.

Stream Language
---------------

This document uses the *Stream Language* to describe messages of
data. This is simple and uniform like JSON, while straightforwardly
supporting strongly typed data like XML.

Documents in the language describe a stream of records. A record
has a name, a set of fields and its own stream of records that are
within it. Fields are always quoted, and always appear within a record, before its record stream.

Some examples:

    empty_record {}
    person { name: "Bob" dob: "1984-16-05" }
    triangle {
        color: "blue"
        point { x: "0" y: "0" }
        point { x: "1" y: "2" }
        point { x: "2" y: "0" }
    }

A BNF:

    stream = record stream 
           |
           ;
    
    record = IDENT '{' fields stream '}'
           ;
    
    fields = IDENT ':' STR fields
           |
           ;




    