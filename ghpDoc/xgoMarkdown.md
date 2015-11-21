<h1 class="appTop">xgoMarkdown</h1>

## What It Does

xgoMarkdown is a tool for generating documents from [xgo templates](template.html).

As the term is used here, a template is a UTF8 document in which the
character sequence `${` (dollar - left bracket) begins a symbol name
and the right bracket `}` ends the symbol name.  Whatever is between
the brackets is first **trimmed** (leading and trailing spaces are removed)
and then used to extract the value of the symbol from the
[context](context.html).

## What It's Used For

## Command Line

...

Explanations for these and other options can be found by typing

    xgoMarkdown -h

and the command line.  Alternatively

    xgoMarkdown -j

will display the current values of command line options and then half.

## Installation
