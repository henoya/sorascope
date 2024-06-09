#!/bin/zsh

HANDLE="henoya.com"
if [[ $# -gt 0 ]]; then
  HANDLE="${1}"
fi

DID="$(bsky show-profile -H "${HANDLE}" --json | jq -r '.did')"
if [[ -z "${DID}" ]]; then
  echo "Handle ${HANDLE} not found"
  exit 1
fi
DISPLAY_NAME="$(bsky show-profile -H "${HANDLE}" --json | jq -r '.displayName' |\
  perl -00 -C -pe 'use utf8;s/[\(\)\[\]\*\N{U+0000}-\N{U+001F}\N{U+007F}-\N{U+009F}\N{U+061C}\N{U+200E}\N{U+200F}\N{U+202A}-\N{U+202E}\N{U+2066}-\N{U+2069}]//g')"

./sorascope -a henoya.com get-all-posts -H "${HANDLE}" --json -n 100
echo "select  created_at, uri, text from post_histories where owner = '${DID}' order by tid desc;" |\
  sqlite3 -json sorascope.db |\
  sed -E 's|("uri":")at://(did:plc:[0-9a-z]+)/app.bsky.feed.post/([0-9a-z]+)"|\1https://bsky.app/profile/\2/post/\3"|'\
    > .result/all-post/${DISPLAY_NAME}_${HANDLE}.all.json
