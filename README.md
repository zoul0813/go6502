# GO6502

Author: David Higgins <https://github.com/zoul0813>

GO6502 is a 6502 Emulator written in Go.
Graphics are provided by [Ebitengine](https://ebitengine.org/).

## Build and run

Requires Go 1.21 or later and the native prerequisites for
[Ebitengine](https://ebitengine.org/en/documents/install.html).

```sh
make build
./bin/go6502
make test
```

The executable includes the original 256-byte Woz Monitor at `$FF00` and the
terminal font. It starts at the ROM reset vector without an assembler, external
ROM file, or a particular working directory. See [ROM provenance](rom/README.md).
The terminal is 40 columns by 24 rows, with RAM at `$0000–$7FFF` and keyboard /
display registers at `$D010–$D013`.

At the monitor prompt, enter `0300: DE AD BE EF`, then `0300.0303` to store and
examine memory. Append `R` to an address to run a program, for example `0300R`.

Home resets the CPU while preserving RAM. Escape restarts monitor input;
Backspace/Delete send the original monitor backspace character. Text is
converted to uppercase ASCII. F1–F4 toggle debug views. F7 toggles single-step mode; F8 executes one instruction
while paused. `--single` starts paused, and `--debug` enables instruction logging.
`--clock 1000` targets one million instructions per second; add `--hz` to express
the rate in instructions per second. Execution and device access share the UI
update loop; timing is not cycle-accurate. Cassette emulation is not implemented.

## Credits

* Graphics Engine: [https://ebitengine.org/](https://ebitengine.org/)
* Font File: [https://style64.org/c64-truetype](https://style64.org/c64-truetype)
* AllSuiteA.asm: [FPGA-netlist-tools](https://github.com/pmonta/FPGA-netlist-tools/blob/master/6502-test-code/AllSuiteA.asm)