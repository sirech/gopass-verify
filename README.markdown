# Gopass Verify

This script is intended to be used to check if your setup will (probably, maybe?) work with [gopass](https://github.com/gopasspw/gopass).

Setting up `gopass` can be a pain in the ass, specially if it is to be used by developers who are not familiar with _gpg_. I have spent a lot of time trying to debug non working setups lately, so I decided to write a script to catch some of the most common problems I have seen.

## Running the script

The script is self documenting, so just run it from the directory

```bash
./go
```

And it will output a list of commands. To run all the possible checks do:

```bash
./go verify
```

The output should hopefully help you figure out what is wrong.

## Fixing errors automatically

This is a much harder problem, and the amount of cases to cover increases a lot. The output of the script should provide some guidance.

## Links

- [Setting up _gopass_](https://hceris.com/storing-passwords-with-gopass/)
