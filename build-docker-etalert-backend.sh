docker build -t etalert-backend:$CI_COMMIT_SHORT_SHA -f  ./Dockerfile .
docker tag etalert-backend:$CI_COMMIT_SHORT_SHA etalert-backend:latest
docker login -u inerxia -p $CI_REGISTRY_PASSWORD
docker image push --all-tags etalert-backend