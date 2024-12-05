# Resource Reduction
You can measure proxy resource consumption (vCPU and memory).

# Preparation
[Common Steps](../docs/common.md)

```bash
./restart.sh
```

# Measurements
```bash
./benchmark.sh <type> <prefix> <rps>
```

# Outputs
You can use `kubectl top` and `k9s` to measure the number of vCPUs and memory.
