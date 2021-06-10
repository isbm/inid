# inid

The "inid" (pronounced as "I need"), is a cross-platform Unix init
system, similar (and inspired by) to `systemd`, written in Go. It can
be used as an extension to `runit` init system (built-in into a
Busybox) or used on its own directly.

The reason to develop this was to have very similar behavior like
systemd with parallel services execution, no need of tracking pids and
forks like in old good SysV. Another goal is to have an init system,
which is using no shell. As well as keep all this really minimal.

## Service Example

Service is just a simple (very-very simple) YAML file with an
extension `.service`. Yet it allows you to compose your own services
out of many commands, not just one, as in `systemd`. Name of a service
is just that: a filename without `.service` extension.

Here is a full example of a service:

```yaml
info: Short description of a service or a title

# Concurrency handling
# after: the service will started after mentioned service done.
#        Omission of "after" means service does not care of dependencies.
# before: same rules as "after".
# NOTE: they are mutually exclusive. If both specified, "before" is ignored.
after: something

# Grouping. You can put unlimited amount of stages. For example,
# assign mounts, agetty, udev etc to the stage 1. Then you can make
# stage 2 for the networking, ssh access, firewalls and so on. The
# "initd" will ensure this sequence.

stage: 1

# Override the default environment.
# The format is "key: value" that will be exported as "key=value".
environment:
  PATH: /sbin:/bin:/usr/sbin:/usr/bin

# Each line in "concurrent" section is a command on its own, running concurrently
concurrent:
  - logger "Hi"
  - logger "something else"

# Each line in "serial" section is a command that is executed after previous
serial:
  - touch /tmp/i-just-booted-it.txt
```

That's it, so far. :wink:
