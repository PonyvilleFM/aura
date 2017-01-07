Aura
====

The stream recording / announcement bot for PonyvilleFM.

Goal
----

Replace Holly as the main method that staff use to record mixes.

### Side goals

- Release this as open source software with clear, clean interfaces into its
  components so that it can be reused by other projects as applicable.
  - There is no advantage to keeping this source code closed-source. The core
    functionality of this is fairly trivial (keep downloading a file until told
    not to, then move it to its final location).
- Create help pages and other user information to make sure that people who
  need to use the recording features know how to use them.

Implementation
--------------

### Information needed

- Twitter API keys [1]
- Stream URL to record
- Web space to put recordings (and HTTPS support for downloads)
- Discord Role ID for users allowed to make recordings (will probably reuse PVFM
  Staff role)

### Commands

#### `.djon [as]`

`.djon` with its optional parameter `as` will start a recording of the currently
live DJ to disk. A recording will by default time out after 4 hours of airtime.
When a name is specified as an argument, that name will become the name that is
announced to twitter as being live.

##### Effects

- Log usage/parameters of the command
- Creates a tweet announcing that the given DJ is online
- Starts a recording process and stores the incoming stream data to the disk

##### Permissions

This command will check its users out by seeing if they are inside a discord role
that allows for recordings to be made or not.

#### `.djoff`

`.djoff` takes no parameters. This command will end a recording with the currently
live DJ, save it, figure out the public-facing URL of the recording, shorten it
and then announce the recording online to Twitter.

##### Effects

- Log usage of the command
- Ends a recording process
- Pushes the recording to public-facing web storage
- Sends out a tweet announcing who was DJing it and the recording URL

##### Permissions

Same as `.djon` permissions.

### State

Persistently, a Recording pointer will be kept. This will be used to trigger
recording stopping when `.djoff` is used.
