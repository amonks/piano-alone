# Controlling the disklavier

This page describes how to download and run the program that talks to the
disklavier. Baseline familiarity with terminal usage is assumed.

## Installation

The disklavier controller program is hosted at
[https://<wbr>piano.computer<wbr>/latest-client](https://piano.computer/latest-client). We can use these commands to download it
into the working directory and make it executable:

```bash
curl https://piano.computer/latest-client -o ./piano-telephone
chmod +x ./piano-telephone
```

## Operation

Then, we can run it:

```bash
./piano-telephone
```

When the program is run, it will connect to the server and check if a
performance is currently taking place. Then, it will present appropriate
options based on that answer. Specifically,

- _If a performance is taking place,_ it will wait for the server to produce the
  final MIDI file, download it, and send it to the disklavier.
- _If a performance is not taking place,_ it will offer to test the MIDI output
  device, and if updates are available, offer to update itself.
