```
 _                   _ _   
| |                 | | |  
| |     ___  _ __ __| | |_
| |    / _ \| '__/ _` | __|
| |___| (_) | | | (_| | |_
\_____/\___/|_|  \__,_|\__|
```

---

Low Orbit Ransomware Denial Toolkit

---

Lordt is an experiment of mine to learn about ransomware and to create anti-ransomware tooling. Consider everything to be "not production ready" for now.

---

# Details

* **Current Version**: 0.0.1
* **License**: MIT
* **Platforms**: Linux
* **Features**
  * `nlock` - An experimental file "lock". It monitors fanotify for file open events and can either simply log when those events happen or also kill the processes performing the `open` syscall on the specified files. Due to fanotify, it isn't fast. A process could read thousands of KB/s before it gets killed. For more information, try `lordt nlock help`.

---

# Usage

## Building

```bash
make build
```

Then the `lordt` executable should be in your current directory.

## `nlock` (experimental, slow)

```bash
./lordt nlock help
```

```bash
sudo ./lordt nlock --pattern test/*.txt -kill -log
```

```bash
sudo ./lordt nlock --conf ./example/nlock-conf.json
```
---



---

This is free and open source software under the MIT license. See `./LICENSE` for more details.