docker build -t inerxia:$CI_COMMIT_SHORT_SHA -f  ./Dockerfile .
docker tag inerxia:$CI_COMMIT_SHORT_SHA inerxia:latest
docker login -u inerxia -p $CI_REGISTRY_PASSWORD
docker image push --all-tags inerxia