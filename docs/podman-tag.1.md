% podman(1) podman-tag - Add tags to an image
% Ryan Cole
# podman-tag "1" "July 2017" "podman"

## NAME
podman tag - Add an additional name to a local image

## SYNOPSIS
**podman [GLOBAL OPTIONS] tag IMAGE[:TAG] TARGET_NAME[:TAG]**
[**--help**|**-h**]

## DESCRIPTION
Assigns a new alias to an image.  An alias refers to the entire image name, including the optional
**TAG** after the ':' If you do not provide a :TAG, podman will assume a :TAG of "latest" for both
the IMAGE and the TARGET_NAME.


## EXAMPLES

  podman tag 0e3bbc2 fedora:latest

  podman tag httpd myregistryhost:5000/fedora/httpd:v2

## SEE ALSO
podman(1), crio(8)

## HISTORY
July 2017, Originally compiled by Ryan Cole <rycole@redhat.com>
