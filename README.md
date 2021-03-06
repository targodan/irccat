# IRCCat

Just a little tool to output things to an IRC server.

# Installing

If you have Go installed simply do `go get -u github.com/targodan/irccat/cmd/irccat`.
If not, well then you'll have to google how to install Go.

# Usage

This will echo the contents of the file called `filename` to the channel `#irccat` on the given server.

```bash
$ irccat server:port filename
```

If you omit the filename it will read from stdin.

Here is a more complete example using TLS and custom channel and username.

```bash
$ echo Hello World. | irccat --tls --nick IRCCatBot --channel '#bots' server:port
```

Note that in order to send the message to a channel you have to include the '#' and pack it in single quotes because otherwise Bash will interpret it as a comment indicator.
This does however mean that you can also send the message to a user instead of a channel.

```bash
$ echo Hello World. | irccat --tls --nick IRCCatBot --channel 'ServerAdmin' server:port
```

# License

MIT License

Copyright (c) 2017 Luca Corbatto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
