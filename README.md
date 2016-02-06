# stalk
A little command-line stalker using the Clearbit API

# Installation

```
go get -u github.com/pzurek/stalk
```

# Config

Sign up with Clearbit and put your key into `~/.stalk/config`:

```
clearbit_key: a78c0257fcfdd142bc1bdda8deadbeef
```

# Use

```
$stalk -e alex@clearbit.com
```

This will hopefully result in something like:

```
Success!
This email seems to belong to: Alex MacCaw
Looks like they are working at Clearbit as a Founder
You can follow them at:
Facebook: https://facebook.com/amaccaw
Twitter:  https://twitter.com/maccaw
GitHub:   https://github.com/maccman
LinkedIn: https://linkedin.com/pub/alex-maccaw/78/929/ab5
```
