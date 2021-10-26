dev:
	go run main.go

test:
	go test ./...

build:
	gcloud builds submit --tag gcr.io/waveblocks/bouncer

deploy:
	gcloud run deploy bouncer \
		--image gcr.io/waveblocks/bouncer \
		--platform managed

ship:
	make test && make build && make deploy
