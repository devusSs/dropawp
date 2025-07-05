# dropawp

A [Go](https://go.dev) CLI tool to track your [Counter-Strike](https://www.counter-strike.net/) inventory prices over time.

## Setup

The usage of this tool is designed to be easy. Download the [latest release](https://github.com/devusSs/dropawp/releases/latest)
and unpack it to a place of your choice.

For some flavors of Linux or also WSL(2) you might need to set up keyring properly first. To do so run each
of the following commands **individually** and **wait for them to complete successfully**.

For the first command **make sure to replace the user id (1000) with your user id**.

```
sudo systemctl restart user@1000.service

sudo apt-get update

sudo apt-get install gnome-keyring libsecret-tools dbus-x11

sudo killall gnome-keyring-daemon

eval "$(printf '\n' | gnome-keyring-daemon --unlock)"

eval "$(printf '\n' | /usr/bin/gnome-keyring-daemon --start)"
```

Thanks to [this little GitHub issue](https://github.com/XeroAPI/xoauth/issues/25) which also helped me resolve those issues.

After that run `dropawp configs add` to create a new tracking project with the corresponding config.
You can simply follow the instructions shown on screen.

## Run

TBA...

## Further usage

Further commands can simply be explored running `dropawp` and checking the printed output. If you have
any questions regarding a subcommand or it's flags you can simply run `<subcommand> help` to get more
output regarding that command, possible further subcommands and flags.

## Disclaimer

This tool is still in active development and does not guarantee bug-free usage. There may also be unintended consequences
or issues created by the tool.
The developer also is not a professional software developer and simply maintains this tool for fun and personal use.
Use this tool at your own responsibility.

Also do not use this tool if you do not understand what it does or which purpose it serves.
Do not use this tool for purposes which might or will violate the [Steam TOS](https://store.steampowered.com/eula/471710_eula_0).

This tool is in no way associated with [Valve](https://www.valvesoftware.com/) or [Counter-Strike](https://www.counter-strike.net/).

## LICENSE

Licensed under [MIT License](./LICENSE).