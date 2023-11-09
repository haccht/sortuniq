# sortuniq

Count the occurence of the lines.
Basically `sortuniq` is equivalent to `sort | uniq` and `sortuniq -cn` is equivalent to `sort | uniq -c | sort -nr`.

```
$ sortuniq -h
Usage:
  sortuniq [OPTIONS]

Application Options:
  -c, --count    Prefix the number of occurrences
  -n, --order    Sort by the number of occurrences
  -r, --reverse  Reverse the order

Help Options:
  -h, --help     Show this help message


$ ps aux | awk '{print $1}' | sortuniq
USER
message+
root
syslog
systemd+
thachimu
www-data

$ ps aux | awk '{print $1}' | sortuniq -cn
  99  root
  14  thachimu
   3  systemd+
   2  www-data
   1  syslog
   1  message+
   1  USER
```

Performance comparison:

```
$ time ( cat /dev/urandom | tr -dc '0-9' | fold -w 5 | head -n 10000000 | sortuniq -cn > /dev/null )

real    0m2.516s
user    0m1.891s
sys     0m0.594s

$ time ( cat /dev/urandom | tr -dc '0-9' | fold -w 5 | head -n 10000000 | sort | uniq -c | sort -nr > /dev/null )

real    0m20.529s
user    0m16.484s
sys     0m1.125s
```
