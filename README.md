# dropawp

A [Go](https://go.dev) CLI tool to track [Counter-Strike](https://www.counter-strike.net/) inventory prices.

It allows you to create a configuration for a tracking project, query your inventory automatically,
add missing items via a file, query prices for each marketable item and store them in a file based database.

It is designed to be easy to use and efficient. You simply create a configuration via the command line,
add your secrets when asked for them and then you can start tracking your inventory prices.

The secrets will be stored securely in [keyring](https://github.com/zalando/go-keyring) (system based password store for Windows, Linux and macOS).

The [Steam API key](https://steamcommunity.com/dev/apikey) is not mandatory in case you do not want to check for required [Steam services](https://steamstat.us/)
before querying your ivnentory. The [CSFloat API key](https://csfloat.com/profile) however is required to query the prices of your items.

## Usage

The usage of this tool is designed to be easy. Download the [latest release](https://github.com/devusSs/dropawp/releases/latest)
and unpack it to a place of your choice.

### Initializing the config and secrets

You may then run `dropawp init` to set up the config and secrets for the tool.

Since the secrets will be managed by keyring this will also need to be set up properly. Usually that works out of the box for Linux, macOS and Windows,
however some flavors of Linux and also WSL(2) have their issues with it. Refer to [the keyring section](./README.md#keyring) for more information.

### Keyring

For some flavors of Linux or also WSL(2) you might need to set up keyring properly first. To do so run each
of the following commands **individually** and **wait for them to complete successfully**.

In case keyring does not work properly for you or e.g. prompts you for a password in a GUI (e.g. on WSL) you can simply run the [fix_keyring.sh script](./scripts/fix_keyring.sh) in the [scripts folder](./scripts/). The easiest way to do this is probably running `/bin/bash scripts/fix_keyring.sh` while being in the repository's directory.

Thanks to [this little GitHub issue](https://github.com/XeroAPI/xoauth/issues/25) which also helped me resolve those issues.

### Running the app

After [setting up](./README.md#initializing-the-config-and-secrets) you can simply run the app using `dropawp run`.

Other implemented commands are available using `dropawp` or `dropawp -h` and may be subject to change in the future. So please refer to the help message for more information.

## Building the app yourself

Although it is highly recommended to download [the latest release](https://github.com/devusSs/dropawp/releases/latest) and simply unpack that to which ever path of your choice and then run the app, you can also build it yourself.

To do so it is highly recommended to run the [included buildscript](./scripts/build.sh) by running `/bin/bash scripts/build.sh` while being in the repository's directory.
That script will set all needed build information if possible and makes the version work / print things properly.

## Disclaimer

This tool is still in active development and does not guarantee bug-free usage. There may also be unintended consequences
or issues created by the tool.
The developer also is not a professional software developer and simply maintains this tool for fun and personal use.
Use this tool at your own responsibility.

This tool is in no way affiliated with [Valve](https://www.valvesoftware.com/) or [Steam](https://store.steampowered.com/) or [CSFloat](https://csfloat.com/).

Make sure to use it with caution and at your own risk. Do not use it for malicious purposes or purposes
that would violate the [Steam TOS](https://store.steampowered.com/eula/471710_eula_0) or [CSFloat TOS](https://csfloat.com/legal/terms-of-service).

## LICENSE

Licensed under [MIT License](./LICENSE).