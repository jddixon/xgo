<h1 class="libTop>Contexts</h1>

An xgo context is a dictionary which maps **symbols** to their **definitions**.
A symbol is a string; a definition may be anything.

Contexts may be nested, so that if a symbol cannot be resolved in 
the immediate context the parent context will be searched and so on
recursively, until there is no parent context.

This also works the other way round: a symbol in a wider context can
be masked by a definition in the local context.

## Operations

Definitions are added to a context at run-time by 

    Context.Bind(symbol string, value interface{}) error

which adds the symbol to the current context, binding it to the 
value specified.

Definitions are retrieved by 

    Context.Lookup(symbol string) (value interface{}, err error)

which returns the definition for the symbol found in the current 
context or should that fail its parent, revursively.

## Serialization

In the current implementation, contexts are stored on disk as simple
name-value pairs: the symbol is followed by a tab which is followed
by the symbol's definition.

