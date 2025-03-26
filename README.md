# Dotsyncer

Sync my dotfiles using git repos

## Commands:

- edit: Edit the config file
- push: Push the changes for each repo

By default it only syncs each folder.

## Config file

The config file is a YAML file stored in $HOME/.config/dotsyncer/config.yaml.

It expects a list with each object having the following attributes:

- name: The name of the config as well the folder name where it will be stored
- repo: The git repo where the files live
- branch: (Optional) the branch to use

### Example Config

```yaml
- name: fish
  repo: git@github.com:kamuridesu/peixefiles.git
  branch: main
- name: nvim
  repo: git@github.com:kamuridesu/nvim-config.git
```
