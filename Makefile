all: website

# This is used to update the external assets.

download-assets: \
	docs/bootstrap.min.css \
	docs/jquery-3.5.1.slim.min.js

docs/bootstrap.min.css:
	curl -L -o $@ https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css

docs/jquery-3.5.1.slim.min.js:
	curl -L -o $@ https://code.jquery.com/jquery-3.5.1.slim.min.js

# Rebuild the index page

index.html: $(wildcard scripts/*)
	cd scripts && yarn && yarn build && yarn release

# Rebuild the visualizer

viz/public/views/index.html:
	cd viz && yarn && yarn build

# Preparing the assets for the website.

define compress =
	gzip -n -c -9 $^ > $@
endef

define copy =
	cp $^ $@
endef

website/docs/%.css: docs/%.css; $(compress)
website/docs/%.js: docs/%.js; $(compress)
website/docs/%.ico: docs/%.ico; $(compress)
website/%.html: %.html; $(compress)

website/docs/%.svg: docs/%.svg; $(copy)
website/docs/%.png: docs/%.png; $(copy)

website/viz/index.html: viz/public/views/index.html; $(compress)
website/viz/public/build/%.js: viz/public/build/%.js; $(compress)
website/viz/public/build/%.css: viz/public/build/%.css; $(compress)

website/viz/public/fonts/%.woff: viz/public/fonts/%.woff; $(copy)

# Actually building the website

website-structure:
	mkdir -p website/docs website/viz/public/build website/viz/public/fonts

website: \
		download-assets \
		website/viz/index.html \
		$(patsubst %,website/%,$(wildcard \
			index.html \
			docs/*.css \
			docs/*.js \
			docs/*.png \
			docs/*.svg \
			docs/*.ico \
			viz/public/build/[!0-9]*.js \
			viz/public/build/*.css \
			viz/public/fonts/*.woff \
		)) \

# Upload website to S3.
#
# Please define the following env var before running:
#
#  - TARGET (in the form `s3://bucket/`)
#  - AWS_PROFILE (the profile name to use, see https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html)
#  - AWS_ENDPOINT_URL (the endpoint URL, default to https://s3.us-east-1.amazonaws.com)

AWS_ENDPOINT_URL ?= https://s3.us-east-1.amazonaws.com

upload: website
	aws --endpoint-url $(AWS_ENDPOINT_URL) \
		s3 sync \
		--exclude '*' \
		--include '*.png' \
		--include '*.svg' \
		--include '*.woff' \
		website/ \
		$(TARGET) && \
	aws --endpoint-url $(AWS_ENDPOINT_URL) \
		s3 sync \
		--exclude '*' \
		--include '*.css' \
		--include '*.js' \
		--include '*.html' \
		--include '*.ico' \
		--content-encoding gzip \
		website/ \
		$(TARGET)

.PHONY: all download-assets website upload

