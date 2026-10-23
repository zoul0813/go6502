# Original Woz Monitor

`monitor/wozmon.rom` is Steve Wozniak's 1976 Apple 1 monitor, exactly 256 bytes,
loaded at `$FF00–$FFFF`. Its reset vector is `$FF00`; the `ECHO` entry is
`$FFEF`. The original NMI and IRQ vectors are preserved, including the
original NMI value `$0F00` (not changed to `$FF00`).

The bytes are decoded from the existing `docs/roms/wozmon.txt` hex listing.
That listing was verified byte-for-byte against `wozmon.txt` in the
[SB-Projects original Apple 1 ROM archive](https://www.sbprojects.net/projects/apple1/roms.zip),
linked from its [downloads page](https://www.sbprojects.net/projects/apple1/download.php).
The [monitor reference](https://www.sbprojects.net/projects/apple1/wozmon.php)
explains the original hardware interface and monitor commands.

`go test ./rom/monitor` verifies the image against the listing and a pinned SHA-256.
The ROM is embedded in the executable; no assembler or ROM download is needed
at build time or runtime.

The older `kernal.rom`, `rom.s`, `wozmon.s`, and `rom/build` are retained as
historical experiments. They are not the runtime ROM or part of the Go build.
