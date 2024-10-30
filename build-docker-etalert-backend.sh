docker build -t docker.io/inerxia/etalert-backend:$CI_COMMIT_SHORT_SHA -f ./Dockerfile .
docker tag docker.io/inerxia/etalert-backend:$CI_COMMIT_SHORT_SHA docker.io/inerxia/etalert-backend:latest
echo $CI_REGISTRY_PASSWORD | docker login -u "inerxia" --password-stdin
docker image push --all-tags docker.io/inerxia/etalert-backend