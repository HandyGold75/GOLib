# Gonos

This module aims to implement an easy and simple way to control Sonos devices while keeping advanced control possible.

The documentation from [svrooij](https://github.com/svrooij/sonos-api-docs) can be referenced for the base implementation.
Many thanks for the documentation as without it this module would not exist.

Next to the base implementation some helper function are present for easier use.
For all helper functions please refer to the helper files.

Note that most of this project is untested.
Some functions might not work as expected.

## Usage

Creating a ZonePlayer for controlling a Sonos device can be done by any of these methods:

```go
zp, err := Gonos.NewZonePlayer("127.0.0.1") // Use the IpAddress of the Sonos device.
zp, err := Gonos.DiscoverZonePlayer() // Discover a Sonos device using SSDP.
zp, err := Gonos.ScanZonePlayer("127.0.0.0/8") // Scan a network for Sonos devices.
```

After a ZonePlayer is successfully created the associated Sonos device can be controlled.

Below a few examples for basic control of a Sonos device:

```go
err := zp.Play() // Play current track.
isPlaying, err := zp.GetPlay() // Check if current track is playing.

err := zp.Pause() // Pause current track.
isPaues, err := zp.GetPause() // Check if current track is paused.

err := zp.Stop() // Stop current track.
isPaues, err := zp.GetStop() // Check if current track is stopped.

isTransitioning, err := zp.GetTransitioning() // Check if current track is transitioning.

err := zp.Next() // Go to next track.
err := zp.Previous() // Go to previous track.

err := zp.SetShuffle(true) // Enable shuffle.
isShuffle, err := zp.GetSuffle() // Check if shuffle is enabled.

err := zp.SetRepeat(true) // Enable repeat (Disables reapeat one).
isRepeat, err := zp.GetRepeat() // Check if repeat is enabled.

err := zp.SetRepeatOne(true) // Enable reapeat one (Disables repeat).
isRepeatOne, err := zp.GetRepeatOne() // Check if repeat one is enabled.

err := zp.SeekTrack(10) // Go to 10th track in the que (Count start at 1).
err := zp.SeekTime(69) // Go to the 69th second in the track.
err := zp.SeekTimeDelta(-15) // Go back 15 seconds in the track.

queInfo, err := zp.GetQue() // Get the current que.
```
