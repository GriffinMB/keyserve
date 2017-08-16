Keyserve serves markdown files from the "blog" directory in Keybase users'
public folder (/keybase/public/{username}/blog). 

The expected directory structure is:

```
blog/
    index.md
    reset.txt
    static/
        style.css
```

To install:

`go get github.com/GriffinMB/keyserve`

To run:

```
$ keyserve -uname=griffinmb -title="Griffin's Blog" &
$ disown %1
```

To clear cache and reset css cache breaker:

`touch blog/reset.txt`
