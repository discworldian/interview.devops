# interview.devops

## Introduction

The assessment is broken down into three sections: [Coding](#coding), [Build](#build), and [Deploy](#deploy). We expect each section to take roughly 40 minutes, for a combined
total of 2 hours. The Coding and Build sections have optional tasks. You can do these as you go if you have time, but consider the tasks ahead and their timing before doing so.

At Popsa we upload a lot of photos. The base code in here defines a very simple Go HTTP server. For the Coding steps, you can adapt the base code to a language of your choosing
but it is encouraged to use Go. Note that the build and deploy steps will continue to use Go, but do not require deep understanding of Go itself. These steps will not depend on the code changes,
so you will be able to do them regardless of your server code.

The next time we meet, we'll run through your solution together. If there are any parts of the assessment that you could not complete within the time, or would have extended / made
more "production ready", please take note of them and we can run through them then.

## Setup

With Go 1.23 installed, the server can be run as such:

```zsh
✗ go run cmd/main.go
2024/09/27 10:52:44 listening on :8081
```

If we run a few requests against the server using the sample images.

```zsh
✗ curl -X POST "localhost:8081/photo?cloudProvider=aws" --data-binary @images/image.jpg
✗ curl -X POST "localhost:8081/photo?cloudProvider=gcp" --data-binary @images/image.png
```

we'll see logs like:

```zsh
2024/09/27 10:53:55 received request with cloudProvider: "aws"
2024/09/27 10:53:55 received request body of size 1097198
2024/09/27 10:54:04 received request with cloudProvider: "gcp"
2024/09/27 10:54:04 received request body of size 605819
```

We can also run it in Docker as such:

```zsh
✗ docker build -t devops-interview .
✗ docker run --rm --name devops-interview -p 8081:8081 devops-interview
```

## Coding

### Task 1

Update the code to return a HTTP Bad Request status code if the `cloudProvider` is not set to one of `aws`, `gcp`, or `azure`.

✅ See photo_upload.go, line 54 onwards (expanded after completing Task 5)

### Task 2

Using a logging framework of your choice, update the logging throughout the code to use a structured log format, e.g. JSON.
If there is a specific reason you have chosen the framework, please make note of that as a comment in the code.

✅ slog seems like the simplest way to go about this. With a larger project (and more time) I would have expanded this out into a helper function.

### Task 3

Update the code to only accept PNG images. If the request body is not a PNG image, return a HTTP Bad Request response.

✅ See photo_upload.go, line 85 onwards.

### Task 4

Add a request query parameter called `bucketName` and using either AWS/GCP/Azure SDK, write code capable of uploading
the photo to the storage bucket.

*Note:* we do not expect you to actually upload the images to a cloud platform. Consider a local solution where possible, for example
[MinIO](https://min.io/product/s3-compatibility) is an S3-compatible storage solution that can be run locally. If you do run such a
solution locally, please include documentation as to how you set it up.

✅ See photo_upload.go, line 85 onwards.

I ran MinIO as suggested, using: 
``` 
docker run
-p 9000:9000
-p 9001:9001
--name minio
-v ~/minio/data:/data
-e "MINIO_ROOT_USER=ROOTNAME"
-e "MINIO_ROOT_PASSWORD=CHANGEME123"
-e "MINIO_SECRET_KEY=CHANGEME123"
-e "MINO_ACCESS_KEY=CHANGEME123"
quay.io/minio/minio server /data 
--console-address ":9001"
```
Which meant I needed to change the docker run command slightly to `docker run --network=host --rm --name devops-interview -p 8081:8081 devops-interview`. 

### Task 5

Update the code so that only one photo upload can be made concurrently to any one cloud provider.

✅ See photo_upload.go, line 54 onwards and line 68/69.

### Task 6 (optional)

Optional task, come back to after build/deploy steps if you have time.

Images may come in rotated, see for example `rotated.png`. Implement code to check if the input image is correctly oriented
and, if not, rotate the image to its default orientation.

❎ I didn't have the time to look into this, but I would look at EXIF orientation packages in go, e.g. https://pkg.go.dev/github.com/takumakei/exif-orientation

### Task 7 (optional)

Optional task, come back to after build/deploy steps if you have time.

Implement unit testing for the the HTTP handler.

❎ I didn't have the time to look at this.

## Build

### Task 1

Despite the end result of our code build being a very small static Go binary, our image is 328MB.

```zsh
✗ docker image ls | grep devops-interview
devops-interview                                                      latest                                     63d1e9df664b   16 minutes ago   328MB
```

Update the Dockerfile to use a multi-stage build to reduce the final binary size.

✅ 

### Task 2

Any change to our Go code, no matter how small, triggers a full rebuild of the Docker image:

```zsh
✗ docker build -t devops-interview .
[+] Building 4.7s (12/12) FINISHED                                                                                                                              docker:desktop-linux
 => [internal] load build definition from Dockerfile                                                                                                                            0.0s
 => => transferring dockerfile: 407B                                                                                                                                            0.0s
 => [internal] load metadata for docker.io/library/golang:1.23-alpine                                                                                                           0.4s
 => [internal] load .dockerignore                                                                                                                                               0.0s
 => => transferring context: 2B                                                                                                                                                 0.0s
 => [1/7] FROM docker.io/library/golang:1.23-alpine@sha256:ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06                                                     0.0s
 => [internal] load build context                                                                                                                                               0.0s
 => => transferring context: 1.11kB                                                                                                                                             0.0s
 => CACHED [2/7] WORKDIR /svc                                                                                                                                                   0.0s
 => [3/7] COPY cmd cmd                                                                                                                                                          0.0s
 => [4/7] COPY internal internal                                                                                                                                                0.0s
 => [5/7] COPY go.* .                                                                                                                                                           0.0s
 => [6/7] RUN go mod download                                                                                                                                                   0.4s
 => [7/7] RUN go build -o /bin/svc cmd/*.go                                                                                                                                     3.6s
 => exporting to image                                                                                                                                                          0.2s
 => => exporting layers                                                                                                                                                         0.1s
 => => writing image sha256:10d85570f6ebe45f70793c39cc466e768c82d16c8fea8af664d6b04d925947a7                                                                                    0.0s
 => => naming to docker.io/library/devops-interview                                                                                                                             0.0s

What's Next?
  View a summary of image vulnerabilities and recommendations → docker scout quickview
✗ sed -i '' -e 's/listen/listening/g' cmd/main.go
✗ docker build -t devops-interview .
[+] Building 4.5s (12/12) FINISHED                                                                                                                              docker:desktop-linux
 => [internal] load build definition from Dockerfile                                                                                                                            0.0s
 => => transferring dockerfile: 407B                                                                                                                                            0.0s
 => [internal] load metadata for docker.io/library/golang:1.23-alpine                                                                                                           0.4s
 => [internal] load .dockerignore                                                                                                                                               0.0s
 => => transferring context: 2B                                                                                                                                                 0.0s
 => [1/7] FROM docker.io/library/golang:1.23-alpine@sha256:ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06                                                     0.0s
 => [internal] load build context                                                                                                                                               0.0s
 => => transferring context: 1.12kB                                                                                                                                             0.0s
 => CACHED [2/7] WORKDIR /svc                                                                                                                                                   0.0s
 => [3/7] COPY cmd cmd                                                                                                                                                          0.0s
 => [4/7] COPY internal internal                                                                                                                                                0.0s
 => [5/7] COPY go.* .                                                                                                                                                           0.0s
 => [6/7] RUN go mod download                                                                                                                                                   0.3s
 => [7/7] RUN go build -o /bin/svc cmd/*.go                                                                                                                                     3.5s
 => exporting to image                                                                                                                                                          0.2s
 => => exporting layers                                                                                                                                                         0.2s
 => => writing image sha256:663ab9b99ebcc5a77e3b702453bd833ff942272e39facb6fa69ca80540cf9e12                                                                                    0.0s
 => => naming to docker.io/library/devops-interview                                                                                                                             0.0s

What's Next?
  View a summary of image vulnerabilities and recommendations → docker scout quickview
```

In a larger project, the time spent in `go mod download` fetching dependencies could be much larger. Update the Dockerfile to cache the dependencies between builds.

✅ 

### Task 3

By default, Docker containers run using the root user. This is explicit in the image `USER root`. Update the image to use a non-root user for improved security.

✅ 

### Task 4

Now that we have our minimal image ready to publish, update the `.github/workflows/publish.yaml` workflow to build and publish the image within the confines of
a Github Actions runner machine, upon merge of a pull request.

*Note:* you don't have to actually publish the image to an image repository.

✅ 

### Task 5 (Optional)

Often, the dependencies that we reference in our Go code at Popsa are private, i.e. they are other Popsa Go projects. This project has a module directive as follows:

```zsh
module github.com/popsa-platform/interview.devops
```

As part of a code change, we have introduced a dependency on a private Go module, e.g. `module github.com/popsa-platform/other-project`. Update the Dockerfile `go mod download` directives
to support downloading private Go modules.

✅

## Deploy

### Task 1

Using your preferred cloud provider among AWS, GCP, or Azure, update the `main.tf` file with resources that can deploy the Docker image as a serverless function.
The Terraform code should consider the following:

* the HTTP server in our code should be reachable at a HTTPS endpoint
* we want to be able to track the cost of running the serverless function

✅

### Task 2

Outline some reasons why we might choose a serverless approach over a long-running service, for example within a Kubernetes service.

✅
