# go-check-longtransaction-cnt

A tool that defines long transactions with a specified threshold and monitors the number of transactions exceeding that value.
Mackerel's exit status can be changed by the number of long transactions.

```
Usage:
  go-check-longtransaction-cnt [OPTIONS]

Application Options:
  -u, --user=       mysql user (default: root)
  -h, --host=       mysql host (default: localhost)
  -p, --port=       mysql port (default: 3306)
      --warn-count= set threshold for warning (default: 0)
      --crit-count= set threshold for critical (default: 0)
      --threshold=  threshold secons
```
