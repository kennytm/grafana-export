# PingCAP MetricsTool

This repository hosts the source code of the website https://metricstool.pingcap.com/, as well as the CLI tool `grafsnap`.

## Deploying the website

You need to first install NodeJS, Yarn 1 and GNU make. The website are pre-compressed, so you need to upload to an S3 server to view it.

```sh
# make sure the Grafana snapshot visualizer is ready.
git submodule update --init

# run `make website` **twice** to build it
make website
make website

# deploy the website on S3
make upload AWS_PROFILE=minio AWS_ENDPOINT_URL='http://127.0.0.1:9000' TARGET='s3://test-website/'
```
