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



## Contributing

Dotstash is, generally speaking, not accepting contributions at this time.
If you think it's cool and want to help out, visit [my website](https://www.jkellogg.dev) and fill out the contact form to drop me a line!
