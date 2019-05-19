# Infrastructure Config Archive v2.2.0

Infrastructure Config Archive (INCA) was developed to solve the problem of
backing up network infrastructure configurations. INCA can be easily expanded to
accommodate multiple types of devices since it uses Expect underneath to handle
the config grabbing.

## Requirements

To Run:

- Expect

To Build:

- Go 1.12+

## Setting INCA

For documentation on setting up INCA, please go to
[http://onesimussystems.com/inca](http://onesimussystems.com/inca).

## Getting Started Developing

```Bash
git clone https://github.com/lfkeitel/inca
cd inca
make build # Build Go application
yarn install # Install and build web frontend
yarn run build
```

## Making a Production Distributable

Clone the repo and run the `package.sh` script. It will build the application
and web frontend and create a compressed tarball in the project root that can
be deployed to a server or Docker container.

## Setup Cron Job

To have configurations pulled on a scheduled basis, you can setup a cron job
that executes:

```Bash
curl http://[hostname]:[port]/api/runnow
```

Set the job to run however often you feel necessary.

## Systemd

There's a baseline Systemd service file in the config directory. You may need
to make edits for settings such as the user/group names, and where the application
is on disk.

## Release Notes

v2.6.0

- Made application paths configurable
- Configuration path can be given as a flag
- Log to stdout as well as files
- Better cross-platform support

v2.5.0

- Better search on archive page
- Consistent JSON from API
- Code restructure for better maintainability

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
