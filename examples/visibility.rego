package opactl.visibility

# Rules whose name starts with __ are invisible on "opactl -a" results.
__config = {
  "artist": "dream theater",
  "songs": ["the alien", "metropolis pt1", "panic attack", "dying soul"]
}

artist = __config.artist

song_id = 1

song = __config.songs[1]

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