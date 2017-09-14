# cda.ms

## Usage

Download a release of the command line tool.

To use it, place the 'cda' or 'cda.exe' binary somewhere in your path.

```
Usage:
  cda URL [flags]
  cda [command]

Available Commands:
  config      create a config file in your home directory
  help        Help about any command
  serve       Run the CDA URL Shortening service

Flags:
  -a, --alias string     CDA Alias
  -c, --channel string   channel
      --config string    config file (default is $HOME/.cda.yaml)
  -e, --event string     event
  -h, --help             help for cda
      --server string    URL Shortening server (default "http://cda.ms")

Use "cda [command] --help" for more information about a command.
```

## Getting Started

Run `cda config`.  This will create a file called `.cda.yaml` in your $HOME directory.

Edit this file, and replace the alias configuration value with your Microsoft alias.

## Creating a Shortened URL

```
cda -c twitter -e ignite https://docs.microsoft.com/azure/x
```

This returns:

```
Using config file: /home/bketelsen/.cda.yaml
Submitting to  http://cda.ms
cda.ms/1
```

It also copies the shortened URL into your clipboard.  Because I love you guys.

## Details

The shortener uses three flags or config settings to create the tracking link.

```
    -a --alias   :  Your CDA alias  ex: brketels
    -c --channel :  The channel/medium  ex: twitter
    -e --event   :  The event name  ex: ignite
```

If the configuration file at `$HOME/.cda.yaml` has any of these values they will be defaulted for you and you may exclude them.

For example, if my config file looks like this:

```
Alias: brketels
Channel: twitter
```
Then I can use the `cda` command and only specify the event with the `-e` flag, the channel and alias will be read from the config file.

You may fill in all three values in the config file, or none.  Any missing values will cause the program to fail with an error:
```
$ cda -a brketels -e ignite https://microsoft.com

> Using config file: /home/bketelsen/.cda.yaml
> Channel is required.  Set with -c or in config file.
```

