

verify: 
	golint
	

test:
	echo "Starting local development server"
	dev_appserver.py app.yaml &


publish: verify
	gcloud app deploy

