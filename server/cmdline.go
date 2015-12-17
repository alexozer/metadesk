package server

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type Cmdline struct {
	root   *Desktop
	server *Server
	args   []string
}

func NewCmdline(root *Desktop, server *Server) *Cmdline {
	return &Cmdline{
		root:   root,
		server: server,
	}
}

func (this *Cmdline) Exec() error {
	if err := this.server.NextConn(); err != nil {
		return err
	}

	if args, err := this.server.ReadCommand(); err != nil {
		return err
	} else {
		this.args = args
	}

	d, err := this.parseDesktop()
	if err != nil {
		return this.server.WriteResponse(err.Error(), ErrorCode)
	}

	return this.parseAndRunCommand(d)
}

func (this *Cmdline) hasArgs() bool {
	return len(this.args) > 0
}

func (this *Cmdline) peek() string {
	return this.args[0]
}

func (this *Cmdline) next() string {
	arg := this.args[0]
	this.args = this.args[1:]
	return arg
}

func (this *Cmdline) checkNoArgs() error {
	if this.hasArgs() {
		var buffer bytes.Buffer

		buffer.WriteString(this.next())
		for this.hasArgs() {
			buffer.WriteString(" ")
			buffer.WriteString(this.next())
		}

		return fmt.Errorf("Warning: extra arguments '%s' ignored", buffer.String())

	} else {
		return nil
	}
}

func (this *Cmdline) parseDesktop() (d *Desktop, err error) {
	if !this.hasArgs() {
		err = errors.New("No desktop selector provided")
		return
	}

	switch this.next() {
	case "root":
		d = this.root
	case "focused":
		d = this.root.Focused()
	case "last":
		d = this.root.LastFocused()
	default:
		err = errors.New("Invalid initial desktop selector")
		return
	}

	for {
		if !this.hasArgs() {
			err = errors.New("No command provided")
			return
		}

		switch this.peek() {
		case "-p", "--parent":
			this.next()

			d = d.Parent()
		case "-c", "--child":
			this.next()

			if !this.hasArgs() {
				err = errors.New("No child index provided")
				return
			}

			childI, e := strconv.Atoi(this.next())
			if e != nil {
				err = errors.New("Invalid child index")
				return
			}

			if !d.IsValidIndex(childI) {
				err = errors.New("Child index out-of-bounds")
				return
			}

			d = d.ChildAt(childI)
		default:
			return
		}

		if d == nil {
			err = errors.New("No matching desktop")
			return
		}
	}
}

func (this *Cmdline) parseAndRunCommand(d *Desktop) error {
	var msg string
	var err error

	switch this.next() {
	case "-f", "--focus":
		d.Focus()

	case "-n", "--next":
		if d.NumChildren() > 0 {
			d.FocusNext()
		} else {
			err = errors.New("Cannot focus next child of leaf")
			break
		}

	case "-N", "--prev":
		if d.NumChildren() > 0 {
			d.FocusPrev()
		} else {
			err = errors.New("Cannot focus previous child of leaf")
			break
		}

	case "-a", "--add":
		d.AddChild()

	case "-r", "--remove":
		if d.NumChildren() > 0 {
			err = errors.New("Cannot remove non-leaf desktop")
		} else if d.IsOccupied() {
			err = errors.New("Cannot remove occupied desktop")
		} else if d == this.root {
			err = errors.New("Cannot remove root desktop")
		} else {
			d.Remove()
		}

	case "-A", "--attrib":
		if !this.hasArgs() {
			err = errors.New("No attribute name provided")
			break
		}
		name := this.next()

		if !this.hasArgs() {
			msg = d.Attr(name)
		} else {
			d.SetAttr(name, this.next())
		}

	case "-u", "--unset":
		if !this.hasArgs() {
			err = errors.New("No attribute name provided")
			break
		}

		d.UnsetAttr(this.next())

	case "-w", "--move-window":
		if d.IsOccupied() {
			d.ClaimFocusedWindow()
		} else {
			err = errors.New("No focused window")
		}

	case "-s", "--swap":
		if d.Parent() == nil {
			err = errors.New("Cannot swap root desktop")
			break
		}

		if !this.hasArgs() {
			err = errors.New("No sibling provided")
			break
		}

		switch sibStr := this.next(); sibStr {
		case "next":
			d.SwapNext()
		case "prev":
			d.SwapPrev()
		default:
			if index, sibErr := strconv.Atoi(sibStr); sibErr == nil {
				if d.Parent().IsValidIndex(index) {
					d.SwapWith(index)
				} else {
					err = errors.New("Invalid sibling index")
				}
			} else {
				err = errors.New("Invalid sibling selector")
			}
		}

	case "-F", "--focused-child":
		msg = fmt.Sprintf("%d", d.FocusedChild())

	case "-C", "--child-count":
		msg = fmt.Sprintf("%d", d.NumChildren())

	case "-P", "--print":
		if !this.hasArgs() {
			err = errors.New("No formatter provided")
			break
		}

		formatter := GetFormatter(this.next())
		if formatter == nil {
			err = errors.New("Unknown formatter")
			break
		}

		msg = formatter.Format(d)

	case "-S", "--subscribe":
		if !this.hasArgs() {
			err = errors.New("No formatter provided")
			break
		}

		formatter := GetFormatter(this.next())
		if formatter == nil {
			err = errors.New("Unknown formatter")
			break
		}

		this.server.SubscribeConn(d, formatter)
		return err

	default:
		err = errors.New("Invalid command")
	}

	if err == nil {
		err = this.checkNoArgs()
	}

	if err == nil {
		return this.server.WriteResponse(msg, SuccessCode)
	} else {
		return this.server.WriteResponse(err.Error(), ErrorCode)
	}
}
