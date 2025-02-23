# Death Star Bench Hotel Reservation
You can compare resource consumption and latency with and without shared microservices using PiCoP.

# Preparation
[Common Steps](../docs/common.md)

# Measurements
Change LUA_NAME in benchmark.sh
## base
```bash
cd ../../DeathStarBench && git checkout base
cd ../benchmark/dsb-hr
./benchmark.sh <PREFIX>
```

## picop
```bash
cd ../../DeathStarBench && git checkout picop
cd ../benchmark/dsb-hr
./benchmark.sh <PREFIX>
```
