<h1 class="appTop">xgoT</h1>

## What It Does

xgoT is a tool for generating documents from [xgo templates](template.html).

As the term is used here, a template is a UTF8 document in which the
character sequence `${` (dollar - left bracket) begins a symbol name
and the right bracket `}` ends the symbol name.  When the context is
**applied**, whatever is between
the brackets is first *trimmed* (leading and trailing spaces are removed)
and then used to extract the value of the symbol from the
[context](context.html).

## What It's Used For

## Command Line

xgoT is a utility for generating documents from templates given a
particular context.  For example,

    xgoT -c x.ctx  -b buildDir a b c

will look for template files a.t, b.t, and c.t in the curent directory,
process them in the context defined by x.ctx, and write the resultant
output files buildDir/a.md, buildDir/b.md, and buildDir/c.m, with
all of the symbolic references, the

    ${variable}

replaced by their definitions.

It is possible to specify

+ the output file extension using **-E ext***, which defaults to `.go`
+ the path to source files using **-t**, which defaults to `./`
+ an output file prefix using **-p**

Explanations for these and other options can be found by typing

    xgoT -h

and the command line.  Alternatively

    xgoT -j

will display the current values of command line options and then half.

## Installation
