# What?

An utility which waits for Internet connection being established.

# Why?

To make systemd units which needs Internet access to feel better.

# How?

The utility pings some public DNS servers and quit after 3 success attempts.
It will fail with exit code 1 if there is no Internet connection established in 5 minutes.

# Build

Go >= 1.20 is required. Run `just`.

# Usage

1. Enable the utility to start: `systemctl enable --now wint`
1. Put these lines into all units which needs an established Internet connection to run:

```
[Unit]
After=wint.service
```

# License

GPL.
