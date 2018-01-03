

verify: 
	golint
	


deploy: verify
	go run *.go


publish: build
	gcloud app deploy

