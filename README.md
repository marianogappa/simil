# simil

String clustering for STDIN

## Usage

```
$ cat file.txt
red
blue
red
blue
red
```

```
$ cat file.txt | simil -k 3
Cluster 0 red
Cluster 0 red
Cluster 0 red
Cluster 1 blue
Cluster 1 blue
```

```
$ cat file.txt | simil -k 3 -short
3 red
2 blue
```

## Syntax

```
simil -k %n [-short] [-random]
```

`-k %n` how many clusters?
`-short` show one line per cluster (shows cluster size + example)
`-random` non-deterministic (randomizes with time-based seed)
