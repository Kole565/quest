BASE64 — КАК ЭТО РАБОТАЕТ

Закодировать:   echo -n "текст" | base64
Раскодировать:  echo "код" | base64 -d

Пример:
  $ echo -n "hello" | base64
  aGVsbG8=
  $ echo "aGVsbG8=" | base64 -d
  hello
