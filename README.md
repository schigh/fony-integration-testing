# fony integration testing

This is a sample project that will show you how to set up a [docker compose](https://docs.docker.com/compose/) environment using [fony](https://github.com/schigh/fony) and [newman](https://www.getpostman.com/docs/v6/postman/collection_runs/command_line_integration_with_newman)



### TL;DR how to

Make sure you have Docker with docker-compose installed.  Clone this repo, then run `docker-compose up` from inside of it.  <kbd>ctrl</kbd>`+`<kbd>c</kbd> to quit.  Deallocate everything by running `docker-compose down`



### what is this?

This repo represents a typical small microservice with a few (5) dependencies.  This service is named `things`, and it has the following dependencies:

```yaml
- caps_svc      # capitalizes things
- foo_svc		# replaces 'foo' with 'bar'
- md5_svc		# generates md5 hash
- reverse_svc	# reverses a string
- sequence_svc	# tells you your index (super helpful)
```

Actual working samples of each service are located under the `externals` folder.  They are mainly just there for reference.

A `fonyfied` service is listed in `docker-compose.yml` as follows:

```yml
  # caps service
  caps_svc:
    image: schigh/fony
    volumes:
      - ./tests/caps-service.json:/fony.json
    ports:
      - "8881:80"
    networks:
      - things_network
```

We supply a suite file for each faked service.  For the above example, we will use this suite file:

```json
{
  "endpoints": [
    {
      "url": "/this is a test",
      "method": "GET",
      "responses": [
        {
          "payload": {
            "value": "THIS IS A TEST"
          }
        }
      ]
    },
    {
      "url": "/of pure foolishness",
      "method": "GET",
      "responses": [
        {
          "payload": {
            "value": "OF PURE FOOLISHNESS"
          }
        }
      ]
    }
  ]
}
```

See [The Suite File](https://github.com/schigh/fony#the-suite-file) for more details.

### newman

We use the [postman/newman](https://hub.docker.com/r/postman/newman/) Docker image for testing within our Docker Compose environment.  To use newman, you must first create a collection in the Postman app, including tests.  More info about Postman tests can be found [here](https://www.getpostman.com/docs/v6/postman/scripts/test_scripts).

Once your collection and tests are ready, you export the collection to a json file.  That is the file that newman runs when the container starts up:

```yaml
  # newman
  newman:
    image: postman/newman
    networks:
      - things_network
    volumes:
      - ./tests/collection.json:/etc/newman/collection.json
    depends_on:
      - things
    command: run /etc/newman/collection.json --bail
```

Here we use our local Postman collection export (`collection.json`) for the test.  The `--bail` flag will force the test to stop as soon as it encounters an error.



Below is the output of a full run of `docker-compose up`:

<details><summary>Contents</summary>
<p>
<pre>
latest: Pulling from schigh/fony
Step 1/11 : FROM golang:alpine AS builder
alpine: Pulling from library/golang
Digest: sha256:692eff58ac23cafc7cb099793feb00406146d187cd3ba0226809317952a9cf37
Status: Downloaded newer image for golang:alpine
 ---> 57915f96905a
Step 2/11 : WORKDIR /
 ---> Using cache
 ---> def1c41fbbe7
Step 3/11 : ADD ./ /
 ---> ae18ee544c00
Step 4/11 : RUN apk add --no-cache git bash
 ---> Running in 1fe2bc7d3e1f
fetch http://dl-cdn.alpinelinux.org/alpine/v3.8/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.8/community/x86_64/APKINDEX.tar.gz
(1/11) Installing ncurses-terminfo-base (6.1_p20180818-r1)
(2/11) Installing ncurses-terminfo (6.1_p20180818-r1)
(3/11) Installing ncurses-libs (6.1_p20180818-r1)
(4/11) Installing readline (7.0.003-r0)
(5/11) Installing bash (4.4.19-r1)
Executing bash-4.4.19-r1.post-install
(6/11) Installing nghttp2-libs (1.32.0-r0)
(7/11) Installing libssh2 (1.8.0-r3)
(8/11) Installing libcurl (7.61.1-r1)
(9/11) Installing expat (2.2.5-r0)
(10/11) Installing pcre2 (10.31-r0)
(11/11) Installing git (2.18.1-r0)
Executing busybox-1.28.4-r1.trigger
OK: 28 MiB in 25 packages
Removing intermediate container 1fe2bc7d3e1f
 ---> 7285543b9dc3
Step 5/11 : RUN mkdir /build
 ---> Running in 3d230e1ba938
Removing intermediate container 3d230e1ba938
 ---> 641017fb8fd2
Step 6/11 : RUN go mod download
 ---> Running in d46a7d36971c
go: finding github.com/schigh/taskchain v1.0.1
go: finding github.com/joho/godotenv v1.3.0
go: finding github.com/pkg/errors v0.8.0
go: finding github.com/ddliu/go-httpclient v0.5.1
go: finding github.com/go-chi/chi v3.3.3+incompatible
Removing intermediate container d46a7d36971c
 ---> c867607de148
Step 7/11 : RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app .
 ---> Running in cbb4fe78910f
Removing intermediate container cbb4fe78910f
 ---> ad52f5aeecfb
Step 8/11 : FROM gcr.io/distroless/base
latest: Pulling from distroless/base
Digest: sha256:8110b784b148db53394d135833ac455c9cc0797e6a143046926de01db3e4c3da
Status: Downloaded newer image for gcr.io/distroless/base:latest
 ---> 48b20972d7e6
Step 9/11 : COPY --from=builder  /build/app /
 ---> 072f455911e0
Step 10/11 : EXPOSE 80
 ---> Running in c62b693b9cb9
Removing intermediate container c62b693b9cb9
 ---> d708429af2ee
Step 11/11 : CMD ["/app"]
 ---> Running in 8e479e32923d
Removing intermediate container 8e479e32923d
 ---> c1521498e92a
Successfully built c1521498e92a
Successfully tagged fony-integration-testing_things:latest
latest: Pulling from postman/newman
Attaching to fony-integration-testing_caps_svc_1_9f3d59faaba8, fony-integration-testing_reverse_svc_1_19658a807e04, fony-integration-testing_md5_svc_1_5c13265d316f, fony-integration-testing_foo_svc_1_b55524243f85, fony-integration-testing_sequence_svc_1_cc6cd5434f56, fony-integration-testing_things_1_e726af8ca0a2, fony-integration-testing_newman_1_fcb9cf20e961
caps_svc_1_9f3d59faaba8 | {"level":"info","ts":"2018-11-28T02:42:50.5906538Z","msg":"fony version: local | build: local"}
caps_svc_1_9f3d59faaba8 | {"level":"info","ts":"2018-11-28T02:42:50.5907322Z","msg":"[fony] - fetching suite from file: fony.json"}
reverse_svc_1_19658a807e04 | {"level":"info","ts":"2018-11-28T02:42:51.2200884Z","msg":"fony version: local | build: local"}
reverse_svc_1_19658a807e04 | {"level":"info","ts":"2018-11-28T02:42:51.2201444Z","msg":"[fony] - fetching suite from file: fony.json"}
md5_svc_1_5c13265d316f | {"level":"info","ts":"2018-11-28T02:42:51.2067864Z","msg":"fony version: local | build: local"}
md5_svc_1_5c13265d316f | {"level":"info","ts":"2018-11-28T02:42:51.206842Z","msg":"[fony] - fetching suite from file: fony.json"}
sequence_svc_1_cc6cd5434f56 | {"level":"info","ts":"2018-11-28T02:42:51.2395879Z","msg":"fony version: local | build: local"}
sequence_svc_1_cc6cd5434f56 | {"level":"info","ts":"2018-11-28T02:42:51.2396765Z","msg":"[fony] - fetching suite from file: fony.json"}
foo_svc_1_b55524243f85 | {"level":"info","ts":"2018-11-28T02:42:51.2268555Z","msg":"fony version: local | build: local"}
foo_svc_1_b55524243f85 | {"level":"info","ts":"2018-11-28T02:42:51.2269342Z","msg":"[fony] - fetching suite from file: fony.json"}
newman_1_fcb9cf20e961 | newman
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | fony-integration-tests
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | → this is a test
md5_svc_1_5c13265d316f | {"level":"info","method":"GET","url":"/this is a test","ts":"2018-11-28T02:42:54.6948283Z","msg":"[fony] - handling request"}
foo_svc_1_b55524243f85 | {"level":"info","method":"GET","url":"/this is a test","ts":"2018-11-28T02:42:54.6964484Z","msg":"[fony] - handling request"}
caps_svc_1_9f3d59faaba8 | {"level":"info","method":"GET","url":"/this is a test","ts":"2018-11-28T02:42:54.6963331Z","msg":"[fony] - handling request"}
reverse_svc_1_19658a807e04 | {"level":"info","method":"GET","url":"/this is a test","ts":"2018-11-28T02:42:54.6971258Z","msg":"[fony] - handling request"}
newman_1_fcb9cf20e961 |   GET http://things/words/this is a test [200 OK, 280B, 124ms]
newman_1_fcb9cf20e961 |   ┌
newman_1_fcb9cf20e961 |   │ '{"caps":{"value":"THIS IS A TEST"},"foo":{"value":"th
newman_1_fcb9cf20e961 |   │ is is a test"},"md5":{"value":"54b0c58c7ce9f2a8b551351
newman_1_fcb9cf20e961 |   │ 102ee0938"},"reverse":{"value":"tset a si siht"}}'
newman_1_fcb9cf20e961 |   │ { caps: { value: 'THIS IS A TEST' },
newman_1_fcb9cf20e961 |   │   foo: { value: 'this is a test' },
newman_1_fcb9cf20e961 |   │   md5:
newman_1_fcb9cf20e961 |   │    { value: '54b0c58c7ce9f2a8b551351102ee0938' },
newman_1_fcb9cf20e961 |   │   reverse: { value: 'tset a si siht' } }
newman_1_fcb9cf20e961 |   └
newman_1_fcb9cf20e961 |   ✓  Response schema is valid
newman_1_fcb9cf20e961 |   ✓  CAPS
newman_1_fcb9cf20e961 |   ✓  FOO
newman_1_fcb9cf20e961 |   ✓  MD5
newman_1_fcb9cf20e961 |   ✓  REVERSE
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | → of pure foolishness
reverse_svc_1_19658a807e04 | {"level":"info","method":"GET","url":"/of pure foolishness","ts":"2018-11-28T02:42:54.8053562Z","msg":"[fony] - handling request"}
foo_svc_1_b55524243f85 | {"level":"info","method":"GET","url":"/of pure foolishness","ts":"2018-11-28T02:42:54.805786Z","msg":"[fony] - handling request"}
md5_svc_1_5c13265d316f | {"level":"info","method":"GET","url":"/of pure foolishness","ts":"2018-11-28T02:42:54.8058258Z","msg":"[fony] - handling request"}
caps_svc_1_9f3d59faaba8 | {"level":"info","method":"GET","url":"/of pure foolishness","ts":"2018-11-28T02:42:54.8078972Z","msg":"[fony] - handling request"}
newman_1_fcb9cf20e961 |   GET http://things/words/of pure foolishness [200 OK, 295B, 9ms]
newman_1_fcb9cf20e961 |   ┌
newman_1_fcb9cf20e961 |   │ '{"caps":{"value":"OF PURE FOOLISHNESS"},"foo":{"value
newman_1_fcb9cf20e961 |   │ ":"of pure barlishness"},"md5":{"value":"4d3271ccb9eec
newman_1_fcb9cf20e961 |   │ 7f7585ca4f6ef16cbc7"},"reverse":{"value":"ssenhsiloof
newman_1_fcb9cf20e961 |   │ erup fo"}}'
newman_1_fcb9cf20e961 |   │ { caps: { value: 'OF PURE FOOLISHNESS' },
newman_1_fcb9cf20e961 |   │   foo: { value: 'of pure barlishness' },
newman_1_fcb9cf20e961 |   │   md5:
newman_1_fcb9cf20e961 |   │    { value: '4d3271ccb9eec7f7585ca4f6ef16cbc7' },
newman_1_fcb9cf20e961 |   │   reverse: { value: 'ssenhsiloof erup fo' } }
newman_1_fcb9cf20e961 |   └
newman_1_fcb9cf20e961 |   ✓  Response schema is valid
newman_1_fcb9cf20e961 |   ✓  CAPS
newman_1_fcb9cf20e961 |   ✓  FOO
newman_1_fcb9cf20e961 |   ✓  MD5
newman_1_fcb9cf20e961 |   ✓  REVERSE
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | → index 0
sequence_svc_1_cc6cd5434f56 | {"level":"info","method":"GET","url":"/{idx}","ts":"2018-11-28T02:42:54.8428406Z","msg":"[fony] - handling request"}
newman_1_fcb9cf20e961 |   GET http://things/sequence/0 [200 OK, 150B, 1014ms]
newman_1_fcb9cf20e961 |   ┌
newman_1_fcb9cf20e961 |   │ '{"value":"your index is 0"}'
newman_1_fcb9cf20e961 |   │ { value: 'your index is 0' }
newman_1_fcb9cf20e961 |   └
newman_1_fcb9cf20e961 |   ✓  Response schema is valid
newman_1_fcb9cf20e961 |   ✓  VALUE
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | → index 1
sequence_svc_1_cc6cd5434f56 | {"level":"info","method":"GET","url":"/{idx}","ts":"2018-11-28T02:42:55.8822724Z","msg":"[fony] - handling request"}
newman_1_fcb9cf20e961 |   GET http://things/sequence/1 [200 OK, 150B, 1509ms]
newman_1_fcb9cf20e961 |   ┌
newman_1_fcb9cf20e961 |   │ '{"value":"your index is 1"}'
newman_1_fcb9cf20e961 |   │ { value: 'your index is 1' }
newman_1_fcb9cf20e961 |   └
newman_1_fcb9cf20e961 |   ✓  Response schema is valid
newman_1_fcb9cf20e961 |   ✓  VALUE
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | → index 2
sequence_svc_1_cc6cd5434f56 | {"level":"info","method":"GET","url":"/{idx}","ts":"2018-11-28T02:42:57.4180874Z","msg":"[fony] - handling request"}
newman_1_fcb9cf20e961 |   GET http://things/sequence/2 [200 OK, 150B, 2s]
newman_1_fcb9cf20e961 |   ┌
newman_1_fcb9cf20e961 |   │ '{"value":"your index is 2"}'
newman_1_fcb9cf20e961 |   │ { value: 'your index is 2' }
newman_1_fcb9cf20e961 |   └
newman_1_fcb9cf20e961 |   ✓  Response schema is valid
newman_1_fcb9cf20e961 |   ✓  VALUE
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | → index 3
sequence_svc_1_cc6cd5434f56 | {"level":"info","method":"GET","url":"/{idx}","ts":"2018-11-28T02:42:59.4540374Z","msg":"[fony] - handling request"}
newman_1_fcb9cf20e961 |   GET http://things/sequence/3 [200 OK, 150B, 2.5s]
newman_1_fcb9cf20e961 |   ┌
newman_1_fcb9cf20e961 |   │ '{"value":"your index is 3"}'
newman_1_fcb9cf20e961 |   │ { value: 'your index is 3' }
newman_1_fcb9cf20e961 |   └
newman_1_fcb9cf20e961 |   ✓  Response schema is valid
newman_1_fcb9cf20e961 |   ✓  VALUE
newman_1_fcb9cf20e961 |
newman_1_fcb9cf20e961 | ┌─────────────────────────┬──────────┬──────────┐
newman_1_fcb9cf20e961 | │                         │ executed │   failed │
newman_1_fcb9cf20e961 | ├─────────────────────────┼──────────┼──────────┤
newman_1_fcb9cf20e961 | │              iterations │        1 │        0 │
newman_1_fcb9cf20e961 | ├─────────────────────────┼──────────┼──────────┤
newman_1_fcb9cf20e961 | │                requests │        6 │        0 │
newman_1_fcb9cf20e961 | ├─────────────────────────┼──────────┼──────────┤
newman_1_fcb9cf20e961 | │            test-scripts │        6 │        0 │
newman_1_fcb9cf20e961 | ├─────────────────────────┼──────────┼──────────┤
newman_1_fcb9cf20e961 | │      prerequest-scripts │        0 │        0 │
newman_1_fcb9cf20e961 | ├─────────────────────────┼──────────┼──────────┤
newman_1_fcb9cf20e961 | │              assertions │       18 │        0 │
newman_1_fcb9cf20e961 | ├─────────────────────────┴──────────┴──────────┤
newman_1_fcb9cf20e961 | │ total run duration: 7.4s                      │
newman_1_fcb9cf20e961 | ├───────────────────────────────────────────────┤
newman_1_fcb9cf20e961 | │ total data received: 435B (approx)            │
newman_1_fcb9cf20e961 | ├───────────────────────────────────────────────┤
newman_1_fcb9cf20e961 | │ average response time: 1195ms                 │
newman_1_fcb9cf20e961 | └───────────────────────────────────────────────┘
fony-integration-testing_newman_1_fcb9cf20e961 exited with code 0
</pre>
</p>
</details>



