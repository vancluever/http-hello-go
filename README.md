# http-hello-go

This is a small web application that I use for demonstration purposes in
various infrastructure projects.

It's the Go equivalent of the demo Sinatra application that I use in my [Packer
and Terraform Example][packer-terraform-example] and [2016 AWS Advent
Demo][advent-demo] repositories, and serves the same purpose as a stand-in
sample web app to demonstrate deployment patterns with [Packer][packer] and
[Terraform][terraform].

[packer-terraform-example]: https://github.com/vancluever/packer-terraform-example
[advent-demo]: https://github.com/vancluever/advent_demo
[packer]: https://www.packer.io/
[terraform]: https://www.terraform.io/

## Building 

A simple `go build ./` will do. You might want to add `CGO_ENABLED=0` to ensure
that it builds static.

You can also control `main.release` to update the version string, which will
influence the contents of the output:

```
go build --ldflags '-X main.release=1.0.0' ./
```

## Testing

There are tests in `main_test.go` as well that can be run with `go test ./`.

## License

```
Copyright 2018-2023 Chris Marchesi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
