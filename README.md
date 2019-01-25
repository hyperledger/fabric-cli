[//]: # (SPDX-License-Identifier: CC-BY-4.0)

**Note:** This is a **read-only mirror** of the formal [Gerrit](https://gerrit.hyperledger.org/r/#/admin/projects/fabric-cli) repository,
where active development is ongoing. Issue tracking is handled in [Jira](https://jira.hyperledger.org/secure/Dashboard.jspa?selectPageId=10104)

This repo is going to be used to implement [FAB-10734 Fabric CLI Redesign](https://jira.hyperledger.org/browse/FAB-10734). This is NOT the "official" Fabric CLI and there is not yet any commitment that it is going to be.

# Hyperledger Fabric CLI

The Hyperledger Fabric CLI is a tool used to interact with [Fabric networks](https://hyperledger-fabric.readthedocs.io/en/latest/).  

## Installation

1. Clone this repo into `$GOPATH/src/github.com/hyperledger`
2. Run `make`
3. Locate the binary in the `bin` directory
4. Add the binary to your PATH
5. Execute `fabric` for more information

## Built-in Commands

Built-in commands can be found in [/cmd/fabric/commands](/cmd/fabric/commands).  These commands can serve as examples for building future commands like `plugin chaincode install ...`.

## Plugins

Users can create and install custom commands to the Fabric CLI.  The only requirement is that all external commands must provide a `plugin.yaml`.

The YAML must specify:
* Name - command name
* Usage - usage syntax
* Description - short description shown for help
* Command - plugin execution

Example plugins can be found in [pkg/plugin/testdata/plugins](pkg/plugin/testdata/plugins).

## Documentation
* [Design Document](https://docs.google.com/document/d/1zIQrS4TRgQEx1z9-wwtO8tYOGRyWdUoTdfk49GFx1wY/edit?usp=sharing)
* [User Stories](https://docs.google.com/document/d/1dxOeM85PgrMNQUJMxB2kwhDthyWnzDxdPvjlwk7x4-w/edit?usp=sharing)

## License <a name="license"></a>

Hyperledger Project source code files are made available under the Apache
License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file.
Hyperledger Project documentation files are made available under the Creative
Commons Attribution 4.0 International License (CC-BY-4.0), available at http://creativecommons.org/licenses/by/4.0/.