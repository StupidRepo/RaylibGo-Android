#!/usr/bin/env bash
set -euo pipefail
adb logcat -d --pid=$(adb shell pidof -s io.github.stupidrepo.RaylibGoGame) -v time
