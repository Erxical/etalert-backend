docker build -t docker.io/inerxia/etalert-backend:$CI_COMMIT_SHORT_SHA -f ./Dockerfile .
docker tag docker.io/inerxia/etalert-backend:$CI_COMMIT_SHORT_SHA docker.io/inerxia/etalert-backend:latest
docker login -u inerxia -p $CI_REGISTRY_PASSWORD
docker image push --all-tags docker.io/inerxia/etalert-backend