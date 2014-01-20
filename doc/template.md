Templates
=========

As the term is used here, a template is a UTF8 document in which the
character sequence ${ (dollar - left bracket) begins a symbol name
and the right bracket } ends the symbol name.  Whatever is between
the brackets is first trimmed (leading and trailing spaces are removed)
and then used to extract the value of the symbol from the context.
The context may be nested, so that if a symbol cannot be resolved in 
the immediate context the parent context will be searched and so on
recursively, until there is no parent context.

In the current implementation there is no provision for escaping
characters between the brackets.  The trimmed, unescaped text between
the brackets must be a valid symbol name.

