# Infrastructure Config Archive v2.2.0

Infrastructure Config Archive (INCA) was developed to solve the problem of
backing up network infrustructure configurations. INCA can be easily expanded to
accommodate multiple types of devices since it uses Expect underneath to handle
the config grabbing.

## Requirements

To Run:

* Expect

To Build:

* Go v1.4

## Setting INCA

For documentation on setting up INCA, please go to
[http://onesimussystems.com/inca](http://onesimussystems.com/inca).

## Getting Started Developing

```Bash
go get github.com/lfkeitel/inca
npm install
```

## Setup Cron Job

To have configurations pulled on a scheduled basis, you can setup a cron job
that executes:

```Bash
curl http://[hostname]:[port]/api/runnow
```

Set the job to run however often you feel necessary. Crontab is the recommended
tool for setting this up and weekly is the recommended schedule.

## Setup Upstart Job

INCA comes with a template upstart script called `upstart.conf`. You can use this
file as a base to build an upstart job to start INCA on boot and to easily manage
the service. Copy the completed script to /etc/init/inca.conf.

You can manage the service with the commands `start inca` `status inca` and `stop
inca`. The upstart job by default will run on boot.

## Release Notes

v2.4.0

- Simpler logging
- Show device name and ip address/hostname when a grab fails
- Keeps an older configuration until a new one doesn't fail

v2.3.0

- Use new Github links
- Notify user if a configuration grab failed

v2.2.0

- Validates format for device configuration files before saving

v2.1.0

- Added application log view on Status page
- Device types can use an asterisk "*" to denote "any method"
- Code cleanup

v2.0.0

- Custom device types: You can now define your own device types and methods.
  Each type/method combo has a script file associated with it that is located
  under the scripts folder. That script will be executed with the arguments
  defined in your device-types.conf.
- Manual device runs: You can manually run a job for a device by entering the
  device information without requiring a full job to be completed first.
- TFTP server no longer required: The Cisco scripts have been changed to no
  longer require a tftp server
- Bug fixes and code cleanup

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

## Versioning

For transparency into the release cycle and in striving to maintain backward
compatibility, This project is maintained under the Semantic Versioning
guidelines. Sometimes we screw up, but we'll adhere to these rules whenever
possible.

Releases will be numbered with the following format:

`<major>.<minor>.<patch>`

And constructed with the following guidelines:

- Breaking backward compatibility **bumps the major** while resetting minor and
  patch
- New additions without breaking backward compatibility **bumps the minor**
  while resetting the patch
- Bug fixes and misc changes **bumps only the patch**

For more information on SemVer, please visit <http://semver.org/>.
