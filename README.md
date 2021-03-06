# go-gcode

go-gcode is an experimental fun project to generate [G-Code](https://reprap.org/wiki/G-code)
for CNC machines and 3D printers, written in Go.

It is loosely based on [Mecode](https://github.com/jminardi/mecode)
and [gcmc](https://www.vagrearg.org/content/gcmc-intro).

Make sure to check out [gocnc](https://github.com/kennylevinsen/gocnc)
which should be able to interpret the G-Code generated by this package.

## Status

This project is just getting started but is usable to generate G-Code
for controlling your CNC machine.

## Usage

The easiest way to use this package is to use a "dot import" to bring
all the GCode functions into your main program's namespace, then
call the functions as if you were drawing each segment.

```go
package main

import (
  "fmt"

  . "github.com/gmlewis/go-gcode/gcode"
)

func main() {
  g := New()
  g.GotoXYZ(XYZ(0,0,0))
  g.MoveX(X(1))
  fmt.Printf("%v", g)
}
```

## Examples

See the [examples](examples) directory for some examples.

To run all the demos at once, type:

```bash
$ ./run-examples.sh
```

----------------------------------------------------------------------

**Enjoy!**

----------------------------------------------------------------------

# License

Copyright 2021 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
