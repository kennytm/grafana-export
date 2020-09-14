#!/bin/sh
set -ex
rm -rf index.html public
mkdir -p public/build public/fonts
cp ../../../3rd_party/grafana/public/views/index.html .
cp ../../../3rd_party/grafana/public/build/*.{js,css} public/build/
cp ../../../3rd_party/grafana/public/fonts/{fontawesome-webfont.woff2,grafana-icons.ttf} public/fonts
