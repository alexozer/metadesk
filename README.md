# metadesk – Organize your Desktops into a Tree

![Metadesk tree](/metadesk.png)

Metadesk is a [bspwm][] wrapper that organizes your desktops into a tree. This is useful for grouping related desktops together to, for instance, separate different tasks. Some features:
- Construct and navigate an arbitrary tree of desktops
- Automatically append and remove empty subdesktops
- Add arbitrary string attributes to desktops, such as names or working directories
- Subscribe to state changes, ready to pipe to status bars like [lemonbar][]

## Installation

You'll first need to set up [go](https://golang.org/doc/install#install) and [bspwm][]. Then run:
```bash
$ go get github.com/alexozer/metadesk/...
```

## Usage

Similar to [bspwm][], metadesk uses a client-server model to control its state. A metadesk daemon, started with the `metadesk` command, is controlled using the `mdc` command. To control metadesk using hotkeys, you use a hotkey daemon like [sxhkd][] to call `mdc`.

#### General `mdc` syntax
```
mdc DESKTOP_SEL COMMAND
```
`mdc` applies the given command to the given selected desktop.

#### Selecting a desktop
A DESKTOP_SEL consists of an initial desktop selection and any number of subselectors:
```
<initial_sel> [SUBSELECTOR]*
```
`<inital_sel>` can be any of:
- `root` the root desktop
- `focused` the currently focused leaf desktop
- `last` the most recently focused leaf desktop

A subselector can be any of:
- `-p | --parent` select the parent of the currently selected desktop
- `-c | --child <n>` select the child desktop at index n

#### Commands
Once a desktop is selected, a command is provided to act on the desktop.

Command | Description
------- | -----------
`-f | --focus` | Focus the selected desktop. If the desktop is not a leaf, focus the last selected child
`-n | --next` | Focus the next child desktop
`-p | --prev` | Focus the previous child desktop
`-a | --add` | Append a child desktop as the last index
`-r | --remove` | Remove the selected desktop
`-A | --attrib <name>` | Print the value of the given attribute
`-A | --attrib <name <value>` | Set the given attribute to the given value
`-u | --unset <name>` | Unset the given attribute
`-w | --move-window` | Move the focused window to the selected desktop
`-s | --swap <sibling>` | Swap the selected desktop with a sibling desktop. `<sibling>` is either `next`, `prev`, or an integer index
`-F | --focused-child` | Print the index of the focused child
`-C | --child-count` | Print the number of child desktops
`-P | --print <formatter>` | Use the given formatter to print the desktop's state
`-S | --subscribe <formatter>` | Subscribe to the desktop's state printed by the given formatter

#### Formatters
Formatters are small functions which format the state of a desktop as a string. The `--print` and `--subscribe` command take a formatter name as an argument. By default, metadesk includes a tree formatter which recursively prints the state of all desktops, and a [lemonbar][] formatter which formats the children of the root desktop for lemonbar. To create your own formatter, you can fork this repository and implement the server.Formatter interface.

#### Attributes
Attributes are user-defined string key-value pairs stored in desktops. They can be accessed by `mdc` and [formatters](#formatters). For example, the default lemonbar formatter prints the `name` attribute of the root desktop's children.

You could use attributes to:
- Store the current working directory of the desktop and open all new terminal windows in this directory
- Store a color string for a formatter to use to output to a panel
- Store some personal notes

## Examples

A sample metadesk configuration which uses [sxhkd][] and [lemonbar][] is available in the [examples](/examples) directory.

#### Example `mdc` usage
Focus the second child of the root desktop:
```bash
mdc root -c 1 -f
```

Select the desktop next to the focused desktop:
```bash
mdc focused -p -n
```

Set the "cwd" attribute of the third child of the root desktop:
```bash
mdc root -c 2 -A cwd "$HOME/code"
```

Remove the last focused desktop:
```bash
mdc last -r
```

Continuously pipe the root desktop's state to lemonbar:
```bash
mdc root -S lemonbar | lemonbar
```

## Contributing

Please do! Pull requests are welcome.

Here are some possible extensions to work on:
- Selectors to select sibling desktops
- Attribute-based desktop selectors
- Multihead support
- React to external desktop focusing
- Additional example formatters

[bspwm]: https://github.com/baskerville/bspwm
[sxhkd]: https://github.com/baskerville/sxhkd
[lemonbar]: https://github.com/LemonBoy/bar
