all: fmt vet

fmt:
	gofmt -s -w .

vet:
	go vet ./...

favicons:
	convert -background none assets/favicons/favicon.svg -resize 16x16 assets/favicons/favicon-16.png
	convert -background none assets/favicons/favicon.svg -resize 32x32 assets/favicons/favicon-32.png
	convert -background none assets/favicons/favicon.svg -resize 180x180 assets/favicons/favicon-180.png
