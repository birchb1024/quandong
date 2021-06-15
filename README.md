# quandong

A program which captures its environment and command-line arguments to a JSON file, then executes the next program in the path with the same name.

## Usage

Copy the quandong executable to a directory using the name of a program whose invocation you wish to capture (the target). Then add the directory to the front of the PATH. Now, whenever `quandong` is executed, its input is captured and the target is loaded with exec() in the same process. Here's usage to log `sudo`:

```
$ mkdir para
$ cp ~/Downloads/quandong para/sudo
$ chmod +x para/sudo
$ export PATH=$PWD/para:$PATH
$ which sudo
/home/lucy/para/sudo
```
When used, there should be no difference from the target.

```
$ sudo id
uid=0(root) gid=0(root) groups=0(root)
$ less quandong-sudo-606199410.json
```

## The Captured Information

Every time the target is run, a temp file is created in the current directory with name format `quandong-<taget name>-<unique number>.json`. For example, `quandong-sudo-606199410.json`

Within the file these attributes are available, example:
```
{
  "args": [         # The command-line arguments
    "sudo",
    "id"
  ],
  "environ": {      # The environment of the process
    "DESKTOP_SESSION": "xfce",
    etc . . .
  },
  "quandong": {     # Information about this program 
    "executable": "/home/lucy/para/sudo",
    "version": "0.0.1-1-g6bd3927"
  },
  "target": "/usr/bin/sudo"     # The path to the target
}

```

## The Name
The quandong is a parasitic tree which draws sustenance from the roots of those around it, like a mistletoe, but it's a standalone tree. It has edible fruit. From Wikipedia:

> Santalum acuminatum, the desert quandong, is a hemiparasitic plant in the sandalwood family, Santalaceae, (Native to Australia) which is widely dispersed throughout the central deserts and southern areas of Australia. 
