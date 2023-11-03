# sortuniq

Count the occurence of the lines.
Basically equivalent to `sort | uniq` or `sort | uniq -c | sort -nr`.

```
$ sortuniq -h
Usage:
  main [OPTIONS]

Application Options:
  -c, --count    Prefix lines by the number of occurrences
  -r, --reverse  Reverse the result

Help Options:
  -h, --help     Show this help message
```

Performance comparison:

```
$ time ( cat /dev/urandom | tr -dc '0-9' | fold -w 5 | head -n 10000000 | sortuniq -c > /dev/null )

real    0m2.516s
user    0m1.891s
sys     0m0.594s

$ time ( cat /dev/urandom | tr -dc '0-9' | fold -w 5 | head -n 10000000 | sort | uniq -c | sort -nr > /dev/null )

real    0m20.529s
user    0m16.484s
sys     0m1.125s
```
