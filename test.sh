# author: https://blog.harrison.dev/2016/06/19/integration-testing-with-docker-compose.html

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

export DISABLE_RATE_LIMITER=true
TIMESTAMP=$(date +%s)

pre_cleanup () { # this made to avoid error with existing containers
  docker rename db db_backup_${TIMESTAMP} || true
  docker rename app app_backup_${TIMESTAMP} || true
}

post_cleanup () {
  docker rename db_backup_${TIMESTAMP} db || true
  docker rename app_backup_${TIMESTAMP} app || true
}

cleanup () {
  docker-compose -p ci kill
  docker-compose -p ci rm -f
  docker network rm enrollment
  post_cleanup
}

trap 'cleanup ; printf "${RED}Tests Failed For Unexpected Reasons${NC}\n"' HUP INT QUIT PIPE TERM

pre_cleanup
docker network create enrollment
docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml build && docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml up -d

if [ $? -ne 0 ] ; then
  printf "${RED}Docker Compose Failed${NC}\n"
  exit -1
fi

TEST_EXIT_CODE=`docker wait ci-tests-1`
docker logs ci-tests-1

if [ -z ${TEST_EXIT_CODE+x} ] || [ "$TEST_EXIT_CODE" -ne 0 ] ; then
  docker logs ci-app-1
  printf "${RED}Tests Failed${NC} - Exit Code: $TEST_EXIT_CODE\n"
else
  printf "${GREEN}Tests Passed${NC}\n"
fi

cleanup
exit $TEST_EXIT_CODE