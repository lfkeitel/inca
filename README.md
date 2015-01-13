Infrastructure Config Archive v1.2.0
====================================

Infrastructure Config Archive (ICA) was developed to solve the problem of backing up network infrustructure configurations.
ICA can be easily expanded to accommodate multiple types of devices since it uses Expect underneath to handle the
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

1. Get the source code (a precompiled build is available at [gzipped archive](/content/ica/infra-config-archive-v1.2.0.tar.gz))
2. Compile with Go
3. Run scripts/setup.sh as root. This will install expect and setup a TFTP server and icauser user account.
4. Copy sample-configuration.toml to configuration.toml
5. Edit the file with the appropiate settings
6. Run executable from directory where you pulled/extracted the application

```Bash
go get github.com/dragonrider23/infrastructure-config-archive
cd $GOPATH/src/github.com/dragonrider23/infrastructure-config-archive
go build
sudo ./scrips/setup.sh
cp sample-configuration.toml configuration.toml
vim configuration.toml
./infrastructure-config-archive
```

Getting Started Developing
--------------------------

```Bash
go get github.com/dragonrider23/infrastructure-config-archive
npm install
```

Setup Cron Job
--------------

To have configurations pulled on a scheduled basis, you can setup a cron job that executes:

```Bash
curl http://[hostname]:[port]/api/runnow
```

Set the job to run however often you feel necessary. Crontab is the recommended tool for setting this
up and weekly is the recommended schedule.

Setup Upstart Job
-----------------

ICA comes with a template upstart script called `upstart.conf`. You can use this file as a base to build an
upstart job to start ICA on boot and to easily manage the service. Copy the completed script to /etc/init/ica.conf.

You can manage the service with the commands `start ica` `status ica` and `stop ica`. The upstart job by default
will run on boot.

Release Notes
-------------

v1.2.0

- Archive a single device configuration
- Added support for Juniper switches
- Bug fixes

v1.1.0

- Edit the device list from the UI
- View application configuration from UI
- Bug fixes

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
