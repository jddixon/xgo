package xmlpull

// Return a copy of the rune slice.

func MakeCopyRunes(src []rune) (dest []rune, err error) {
	if src == nil || len(src) == 0 {
		err = EmptyRuneSlice
	} else {
		dest = make([]rune, len(src))
		copy (dest, src)
	}
	return
}

// Return a copy of the tag name (in which case the argument is nil 
// or empty, or copy the parameter as an entity ref name.
	
func (p *Parser) getName(candidate []rune) (runes []rune, err error) {

    if p.curEvent == START_TAG {
        runes, err = MakeCopyRunes( p.elName[ p.elmDepth ] )
    } else if p.curEvent == END_TAG {
        runes, err = MakeCopyRunes( p.elName[ p.elmDepth ] )
    } else if p.curEvent == ENTITY_REF {
        if p.entityRefName == nil {
			runes, err = MakeCopyRunes(candidate)
			if err == nil {
				p.entityRefName, err = MakeCopyRunes(runes)
			}
        }
    }
	if err != nil {
		runes = nil
	}
	return
}

func (p *Parser) getNamespace() (runes []rune, err error) {

    if p.curEvent == START_TAG || p.curEvent == END_TAG {
        if p.processNamespaces {
			runes, err = MakeCopyRunes(p.elUri[ p.elmDepth  ])
		} 
	}
	if runes == nil {
		runes, err = MakeCopyRunes(NO_NAMESPACE)
	}
    return 
} 
