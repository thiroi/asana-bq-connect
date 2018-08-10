APP_NAME=goappengine-sample

# install modules (≒ npm install)
install:
	git submodule update --init --recursive

# run local
run:
	goapp serve --host=0.0.0.0 --port=8080 ./

# run local with clear database
clean-run:
	goapp serve --clear_datastore --host=0.0.0.0 --port=8080 ./

# deploy to appengine
deploy:
	goapp deploy -application ${APP_NAME} ./
