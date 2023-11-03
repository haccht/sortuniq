# sortuniq

Count the occurence of the lines.
Basically equivalent to `sort | uniq` or `sort | uniq -c | sort -n`.

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
