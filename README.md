A simplistic live streaming origin written in Go

- Upload `PUT` to `/ingest/`
- Access `GET` from `/live/`

## Example

```
curl -v -X PUT -H "Content-Type: application/octet-stream" --data-binary '@testmedia/live/test/2M/fileSequence422.mp4' http://localhost:8080/ingest/test2/2M/fileSequence422.mp4

curl -v http://localhost:8080/live/test2/2M/fileSequence422.mp4
```

## Build and Run

```
$ go build

$ ./go-origin
2020/10/21 14:29:32 Starting Eyevinn simple origin store=./testmedia/
```

Set environment variable `MEDIAPATH` if you don't want to use the default `./testmedia/`