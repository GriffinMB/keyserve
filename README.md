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

Debugging note:

If you are having problems getting Keybase to start on Ubuntu, you may need to create the following file:

```
$ touch /home/<user>/.cache/keybase/keybase.kbfs.log
```
