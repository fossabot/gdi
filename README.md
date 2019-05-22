GoDropIt
========

[![Build Status](https://gitlab.com/ilpianista/gdi/badges/master/build.svg)](https://gitlab.com/ilpianista/gdi/pipelines)

A command-line client for [DropIt](https://apps.nextcloud.com/apps/dropit) Nextcloud app.

## Get it

    $ go get gitlab.com/ilpianista/gdi
    $ go install gitlab.com/ilpianista/gdi

## Usage

    $ echo "Hey, look at GoDropIt" | gdi -s https://mynextcloud -u user
    Generating link... https://mynextcloud/s/of6tzHPwCHsc6Po

or

    $ gdi -s https://mynextcloud -u user -f /path/to/mybinaryfile
    Generating link... https://mynextcloud/s/of6tzHPwCHsc6Po

## Donate

Donations via [Liberapay](https://liberapay.com/ilpianista) or Bitcoin (1Ph3hFEoQaD4PK6MhL3kBNNh9FZFBfisEH) are always welcomed, _thank you_!

## License

MIT
