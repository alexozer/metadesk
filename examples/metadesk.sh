#!/bin/bash

mdc root -a
mdc root -a
mdc root -a
mdc root -a

# add a child desktop to prevent automatic removal
mdc root -c 0 -a
mdc root -c 1 -a
mdc root -c 2 -a
mdc root -c 3 -a

# set names of parent desktops
mdc root -c 0 -A name email
mdc root -c 1 -A name chat
mdc root -c 2 -A name web
mdc root -c 3 -A name code

# show clickable parent desktop names in lemonbar
mdc root -S lemonbar |\
	lemonbar -g 1000x30+32 -B '#2D2D2D' -F '#EDEDED' -f 'sans:size=10' |\
	while read index; do
		mdc root -c $index -f
	done

