A simplistic live streaming origin written in Go

- Upload `PUT` to `/ingest/`
- Access `GET` from `/live/`

## Example

Upload:

```
curl -v -X PUT -H "Content-Type: application/octet-stream" --data-binary '@testmedia/live/test/2M/fileSequence422.mp4' http://localhost:8080/ingest/test2/2M/fileSequence422.mp4
```

Access:

```
curl -v http://localhost:8080/live/test2/2M/fileSequence422.mp4
```

Delete:

```
curl -v -X DELETE http://localhost:8080/ingest/test2/2M/fileSequence422.mp4
```

## Build and Run

```
$ go build

$ ./go-origin
2020/10/21 14:29:32 Starting Eyevinn simple origin store=./testmedia/
```

Set environment variable `MEDIAPATH` if you don't want to use the default `./testmedia/`.
Set environment variable `PORT` to specify what port to listen to. Default is 8080.

## About Eyevinn Technology

Eyevinn Technology is an independent consultant firm specialized in video and streaming. Independent in a way that we are not commercially tied to any platform or technology vendor.

At Eyevinn, every software developer consultant has a dedicated budget reserved for open source development and contribution to the open source community. This give us room for innovation, team building and personal competence development. And also gives us as a company a way to contribute back to the open source community.

Want to know more about Eyevinn and how it is to work here. Contact us at work@eyevinn.se!
