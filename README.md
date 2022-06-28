# vanity-keygen

Finds a vanity SSH ed25519 keypair. If you want your public key to contain something like `pCaRR1er`, run with:

```
$ go install github.com/pcarrier/vanity-keygen@latest
$ ~/go/bin/vanity-keygen '[pP][cC][aA][rR][rR][iI1][eE][rR]'
```

and wait for the program to terminate, depending on your luck
after tens of billions (>1e10) of generations in this instance,
fewer than a hundred millions (<1e8) for `[Pp]ierre`,
millions (1e7) for `love$`,
fewer than a million (<1e6) for `love`.

Use `-threads N` if you don't want to saturate all cores.

On servers, `tmux` is your friend. `C-b z` to zoom for easy copy-paste.

Assuming you saved the private key in `vanity`, remember to set a passphrase with:

```
$ ssh-keygen -p -f vanity
```

Alternatively, use the comment (third) field in your `.pub` and `authorized_keys`. But where's the fun in that?

## Security

This is reminiscent of bitcoin mining.
Finding a SHA-256 hash with a lot of zeroes through trial and error doesn't make it any more prone to collisions.
We are similarly using brute force to find a public key of a particular shape,
which shouldn't make the private key any more discoverable nor its security any weaker.

Do worry about a potential capture of the program's output, and cleartext storage of the private key.

As such, don't run this on a multi-user system for a personal keypair.

At the risk of repeating myself, do not use persistent storage that's not encrypted or physically secured
throughout its lifetime to buffer the program's output or store the plaintext private key (however temporarily).

Of course, don't blindly trust random strangers on the Internet.
My chunk of the code, [`main.go`](main.go), is tiny and should make for a quick read;
`go mod verify` to ascertain `vendor/` isn't compromised, `go mod graph` to explore what's there.

Not much else for me to say about your Go standard library, runtime, toolchain, OS, and hardware.
Open is better. [Security is hard.](https://wiki.c2.com/?TheKenThompsonHack)

## Performance

According to `perf top -g` and Instruments.app, we're spending the majority of the time in ed25519 arithmetic,
so the rest (string manipulations, allocations, etc.) doesn't seem worth much effort.

## Example output

```
$ vanity-keygen 'love$'
2022/07/01 03:45:24 Looking for a public key matching love$
2022/07/01 03:45:25 Generated 232,000 keypairs
[…]
2022/07/01 03:48:00 Generated 36,335,000 keypairs
2022/07/01 03:48:01 Public key:
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBghf18a0v58TKzdx7urZ1Lv3cfQfySi8WCHaSdClove
2022/07/01 03:48:01 Private key:
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtz
c2gtZWQyNTUxOQAAACAYIX9fGtL+fEys3ce7q2dS793H0H8kovFgh2knQpaL3gAA
AIiaywRCmssEQgAAAAtzc2gtZWQyNTUxOQAAACAYIX9fGtL+fEys3ce7q2dS793H
0H8kovFgh2knQpaL3gAAAEAGCrt1RnJFV4IigAwebCePX4wDo1HNnQGsG5tV8JAY
aBghf18a0v58TKzdx7urZ1Lv3cfQfySi8WCHaSdCloveAAAAAAECAwQF
-----END OPENSSH PRIVATE KEY-----
```
