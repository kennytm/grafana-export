all: download-assets

download-assets: \
	docs/bootstrap.min.css \
	docs/jquery-3.5.1.slim.min.js

docs/bootstrap.min.css:
	curl -L -o $@ https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css

docs/jquery-3.5.1.slim.min.js:
	curl -L -o $@ https://code.jquery.com/jquery-3.5.1.slim.min.js

.PHONY: all download-assets

