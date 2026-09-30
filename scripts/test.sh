#!/bin/bash

FAILED=0

run_test() {
  local description="$1"
  shift
  echo "Running Test: $description"
  if "$@"; then
    echo "  ✅ Passed"
  else
    echo "  ❌ FAILED"
    FAILED=$((FAILED + 1))
  fi
}

can_agent_start_nginx() {
  sudo -u agent sudo systemctl start nginx && sudo -u agent sudo systemctl stop nginx
}

locked_archive_exist() {
  test -f /srv/config/wheel.tar.gz && tar -tf /srv/config/wheel.tar.gz > /dev/null 2>&1
}

echo "Running Testing suite for a Quest"

run_test "does readme exist" test -s /home/agent/readme
run_test "is nginx installed" test -f /etc/nginx/nginx.conf
run_test "is nginx shutted down" bash -c '! systemctl is-active --quiet nginx'
run_test "is server can be start from agent" can_agent_start_nginx
run_test "index page exist" test -s /var/www/html/index.html
run_test "locked .tar archive exist" locked_archive_exist
run_test "is java installed" command -v java > /dev/null 2>&1
run_test "does mc alias exist" grep -qE "^\s*alias\s+ddd\s*=\s*'java /opt/ddd/ddd.jar'" /home/agent/.bashrc

if [ "$FAILED" -eq 0 ]; then
  echo "Test suite is done ✅"
else
  echo "Test suite is completed with errors ❌"
fi

# exit $FAILED

