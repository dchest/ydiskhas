#!/bin/sh

which curl >/dev/null || { echo "curl not found" >&2; exit 126; }
AUTH=$2; [ -z $AUTH ] && { echo "Please specify login or login:password as second argument" >&2; exit 10; }

FILE=$1; SIZE=`stat -c%s "$FILE" 2>/dev/null || stat -f%z "$FILE" 2>/dev/null` || { echo "Could not determine size of $FILE!" >&2; exit 11; }

MD5=`{ md5 "$FILE" 2>/dev/null || md5sum --tag "$FILE"; } | sed 's/.* //'`
[ -z $MD5 ] && { echo "Cannot get file MD5" >&2; exit 12; }
SHA256=`shasum -a 256 $FILE | sed 's/ .*$//'`;
[ -z $SHA256 ] && { echo "Cannot get file SHA256" >&2; exit 12; }

URL="https://webdav.yandex.ru/%3Fexists%3F$SHA256"
STATUS=`echo -n | curl -T- -u "$AUTH" -X PUT $URL -H "Content-Range: bytes 0-0/$SIZE" -H Etag:$MD5 -H Sha256:$SHA256 -s -w "\n%{http_code}\n" | tail -1`

test $STATUS = 201 || { echo "(╯°□°)╯︵ ┻━┻\nFile does not exist"; exit 1; }
echo "˙ ͜ʟ˙\nFile exists. Removing..."
curl -u "$AUTH" -X DELETE $URL; exit 0
