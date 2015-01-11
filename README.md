Infrastructure Config Archive v1.0.0
====================================

Infrastructure Config Archive (ICA) was developed to solve the problem of backup network infrustructure configurations.
ICA can be easily expanded to accomidate multiple types of devices since it uses Expect underneath to handle the
config grabbing.

Requirements
------------

To Run:

* Expect
* tftp server

To Build:

* Go v1.4

Is it any good?
---------------

[Yes](https://news.ycombinator.com/item?id=3067434)

Setting Up and Using ICA
------------------------

1. Get the source code (current there's no precompiled binaries)
2. Compile with Go
3. Copy sample-configuration.toml to configuration.toml
4. Edit the file with the appropiate settings
5. Run executable from directory where you pulled/extracted the application

```Bash
go get github.com/dragonrider23/infrastructure-config-archive
cd $GOPATH/src/github.com/dragonrider23/infrastructure-config-archive
go build
cp sample-configuration.toml configuration.toml
vim configuration.toml
./infrastructure-config-archive
```

Getting Started Developing
--------------------------

```Bash
go get github.com/dragonrider23/infrastructure-config-archive
npm install
go get
```

Setup Cron Job
--------------

To have configurations pulled on a schedules basis, you can setup a cron job that executes:

```Bash
curl http://[hostname]/api/runnow
```

Set the job to run however often you feel necassary. Crontab is the recommend tool for setting this
up and weekly is the recommeneded schedule.

Release Notes
-------------

v1.0.0

- Initial Release

Versioning
----------

For transparency into the release cycle and in striving to maintain backward compatibility, This project is maintained under the Semantic Versioning guidelines. Sometimes we screw up, but we'll adhere to these rules whenever possible.

Releases will be numbered with the following format:

`<major>.<minor>.<patch>`

And constructed with the following guidelines:

- Breaking backward compatibility **bumps the major** while resetting minor and patch
- New additions without breaking backward compatibility **bumps the minor** while resetting the patch
- Bug fixes and misc changes **bumps only the patch**

For more information on SemVer, please visit <http://semver.org/>.
