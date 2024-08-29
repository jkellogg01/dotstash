# Dotstash!

Your friendly configuration gardener.

Your dotfiles are precious.
You modify them carefully, you share them joyfully.
In return, they provide your computing experience with the _je ne sais quois_ that puts a smile on your face every time you log in.

Dotstash is a program designed to manage all of the busywork associated with keeping your dotfiles safe.
Go ahead, configure to your heart's content.
We'll take care of the rest.

## Disclaimer

Dotstash is in _very early development_.
The final product should be incredibly hands-off and useful even to those without much experience using CLI utilities.
That being said, this is not the final product.
The last thing I want is for anyone to lose data on my account, so here are some recommendations while the kinks are still being worked out:

### Backups

Make copies of anything you are going to put into dotstash, **before** you do so!
You should probably be keeping more robust backups than that anyways, but I certainly don't and I'm not your dad.

### CLI Comfort

If you are not highly comfortable with command-line utilities or using the terminal to navigate your filesystem, Dotstash may be a bit out of your depth.
There are plans to fix this with an interactive TUI in the near future, but in the meantime trying to use Dotstash may do more harm than good.

### TL;DR

Don't be stupid.
Don't trust a stranger's software with data you can't afford to lose or accidentally share.
That's about it.

GLHF! <3

## A Note on Terminology

Throughout this document as well as the help menu for many of the commands, you will find a bountiful harvest of gardening terminology.
This is by design, as many elements of your Dotstash configuration will be strikingly similar to concepts you're already familiar with if you're sufficiently comfortable with git and with navigating filesystems, especially in the terminal.
As such, gardening terminology has been borrowed to clearly distinguish the nuanced elements of Dotstash from their native analogues.

- **Gardens** encompass a _collection_ of configuration files. A garden is confined to a single branch of a single git repository. Dotstash provides the tools to interact with the underlying git repository from anywhere, so you don't even need to go find where we keep them (they should be in a directory called `.dotstash` inside of your home directory)!
- **Flowers** refer to the set of configuration files for _a single application_. It may just be `~/.zshrc`, or it may be your whole `~/.config/nvim` directory; both of these are **Flowers** as far as Dotstash is concerned.

Understanding this will make commands such as `dotstash plant ~/.config/wezterm` or `dotstash uproot neofetch` (forever in our hearts) entirely sensible and easy to distinguish from typical filesystem commands.

## Installation

Dotstash can be installed using go:

```sh
go install github.com/jkellogg01/dotstash@latest
```

More, easier ways to install are coming in the future, especially as Dotstash nears v1.0!

## Usage

Dotstash's usage is made up of a relatively small list of commands:

### Git

The `git` command executes arbitrary git commands inside of the primary garden.
The `garden` flag **must be immediately after `git` and before the git command itself**, and allows you to specify an alternate garden in which to execute the git command.
In order to commit some changes to my configuration in a garden called `franklin` (idk man, a name is a name), I would use the following commands:

```sh
dotstash git --garden=franklin add -A
dotstash git --garden=franklin commit
```

Note that `git` is compatible with interactive git commands, so it will have no problem using my default editor to generate a commit message.

### Make (Sow)

The `make` command creates a new garden with the flowers specified by the arguments.
The `name` flag allows you to specify a name for the garden. The name will default to `dotstash`.
The `author` flag allows you to specify an author for the garden. The author will default to your username.
For example, if I wanted to create a garden called `nv` which contained my neovim configuration, I would do so with the following command:

```sh
dotstash make ~/.config/nvim --name=nv
```

### Plant

The `plant` command adds a flower to the primary garden.
The `garden` flag allows you to specify an alternate garden in which to plant the flower.
For example, if I wanted to add my `.zshrc` to the `nv` garden we created in the previous example, I would use the following command to do so:

```sh
dotstash plant ~/.zshrc --garden=nv
```

Note that if `nv` was my primary garden I would not need to specify it with the `garden` flag.

### List

The `list` command prints a list of your current gardens and all of the flowers listed therein.

### Uproot (Deplant)

The `uproot` command removes a flower from the primary garden.
The `garden` flag allows you to specify an alternate garden from which to remove the flower.
For example, if I wanted to remove my `.zshrc` from the `nv` garden, I would execute the following command:

```sh
dotstash uproot .zshrc --garden=nv
```

### Select

The `select` command selects a garden to use as a source for your system's configuration style.
The `clobber` flag will write over non-flower configuration files contained by the newly-selected garden.
The `unlink` flag will preserve flowers contained by the newly-selected garden, _instead of_ unlinking them.
In order to set the `nv` garden from the previous examples as primary, I would use the following command:

```sh
dotstash select nv
```

### Remove

The `remove` command removes a garden from your list of gardens.
The `no-restore` flag causes the configuration files to be deleted permanently.
In order to remove the `nv` garden from the previous examples, I would use the following command:

```sh
dotstash remove nv
```

### Torch

Do not use the `torch` command.
It is a destructive tool created for my personal convenience while testing.
You have been warned.

## Contributing

Dotstash is, generally speaking, not accepting contributions at this time.
If you think it's cool and want to help out, visit [my website](https://www.jkellogg.dev) and fill out the contact form to drop me a line!
