package opactl.examples.visibility

# A comment for subcommand "visibility" 
__comment = "Example of how to use invisible field __ and comments"

# Rules whose name starts with __ are invisible on "opactl -a" results.
__config = {
  "artist": "dream theater",
  "songs": ["the alien", "metropolis pt1", "panic attack", "dying soul"]
}

# A comment for subsubcommand "artist"
__artist = "Get data from __config field"
artist = __config.artist

# A comment for subsubcommand "song_id"
__song_id = "Get id of song"
song_id = 1

# A comment for subsubcommand "song"
__song = "Get a song specified with song_id"
song = __config.songs[song_id]

## Output examples: 
# $ opactl visibility -a
#[
#  "artist",
#  "song",
#  "song_id"
#]

# $ opactl visibility artist
# "dream theater"

# $ opactl visibility song
# "metropolis pt1"