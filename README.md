# console-notifications
Print desktop notifications in a console

```
$ notify-send --app-name test-app1 "summary text" "body text\n\nmore body text"
$ notify-send --app-name test-app2 "summary text2" "body text 2"
```

```
Monitoring notifications. Enter to clear screen. Ctrl-c to quit.
-------------------------------------------------------------------------------------------
2023-03-01 16:02:33 - test-app1

summary text

body text

more body text

-------------------------------------------------------------------------------------------
2023-03-01 16:03:13 - test-app2

summary text2

body text 2

-------------------------------------------------------------------------------------------

```
