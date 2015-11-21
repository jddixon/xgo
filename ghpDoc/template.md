<h1 class="libTop">Templates</h1>

As the term is used here, a template is a UTF8 document in which the
character sequence `${` (dollar - left bracket) begins a symbol name
and the right bracket `}` ends the symbol name.  Whatever is between
the brackets is first **trimmed** (leading and trailing spaces are removed)
and then used to extract the value of the symbol from the
[context](context.html).

## Limitations

In the current implementation there is no provision for escaping
characters between the brackets.  The trimmed, unescaped text between
the brackets must be a valid symbol name.

## Command Line Tools

* [xgoT](xgoT.html)


