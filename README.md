# dotstash!

a go program for managing dotfiles conveniently

## TODO

- [x] make command to create a new directory
- [x] make command should move targeted config files to new directory and symlink them back to their original locations
- [x] make command should initialize a git repository in the created directory
- [x] remove command to delete a config repository and, if necessary, move source files back to their 'targets'
- [x] depend[^1] command to add a config file or directory to a repository
- [x] obviate[^1] command to remove a config file or directory from a repository
- [x] list command to list all config repositories
- [x] select command to select a 'primary' config and link the config files from it
- [ ] get command to download a remote repository and, if possible, resolve targets
- [ ] style/terminology pass:
  - [ ] strong definition of commonly-used terminology
  - [ ] more consistent use of a more specific set of APIs and standard library packages
- [ ] more visual polish; interactive modes for commands that didn't need them before

[^1]: the terms 'depend' and 'obviate' might be a bit obscuring here, but I wanted them to be substantially different from 'make' and 'remove' and so those are the terms I'm using.
