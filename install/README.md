# Istio installation

This directory contains the default Istio installation configuration and
the script for updating it.
 
## updateVersion.sh

The [updateVersion.sh](../updateVersion.sh) script is used to update image versions in
[../istio.VERSION](../istio.VERSION) and the istio installation yaml files.

### Options

* `-m <hub>,<tag>` new manager image
* `-x <hub>,<tag>` new mixer image
* `-i <url>` new `istioctl` download URL
* `-c` create a `git commit` titled "Updating istio version" for the changes

Default values for the `-m`, `-x`, and `-i` options are as specified in `istio.VERSION`
(i.e., they are left unchanged).

### Examples

Update the manager and istioctl:

```
./updateVersion.sh -m "docker.io/istio,2017-05-09-06.14.22" -c "https://storage.googleapis.com/istio-artifacts/dbcc933388561cdf06cbe6d6e1076b410e4433e0/artifacts/istioctl"
```
