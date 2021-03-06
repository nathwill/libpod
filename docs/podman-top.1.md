% podman(1) podman-top - display the running processes of a container
% Brent Baude

## NAME
podman top - Display the running processes of a container

## SYNOPSIS
**podman top**
[**--help**|**-h**]

## DESCRIPTION
Display the running process of the container. ps-OPTION can be any of the options you would pass to a Linux ps command

**podman [GLOBAL OPTIONS] top [OPTIONS]**

## OPTIONS

**--help, -h**
  Print usage statement

**--latest, -l**
Instead of providing the container name or ID, use the last created container. If you use methods other than Podman
to run containers such as CRI-O, the last started container could be from either of those methods.

## EXAMPLES

```
# podman top f5a62a71b07
  UID   PID  PPID %CPU STIME TT           TIME CMD
    0 18715 18705  0.0 10:35 pts/0    00:00:00 /bin/bash
    0 18741 18715  0.0 10:35 pts/0    00:00:00 vi
#
```

```
#podman --log-level=debug top f5a62a71b07 -o pid,fuser,f,comm,label
  PID FUSER    F COMMAND         LABEL
18715 root     4 bash            system_u:system_r:container_t:s0:c429,c1016
18741 root     0 vi              system_u:system_r:container_t:s0:c429,c1016
#
```
## SEE ALSO
podman(1), ps(1)

## HISTORY
December 2017, Originally compiled by Brent Baude<bbaude@redhat.com>
