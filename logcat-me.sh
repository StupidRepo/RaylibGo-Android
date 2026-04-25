#!/usr/bin/env bash
set -euo pipefail

clear
adb logcat -d --pid=$(adb shell pidof -s io.github.stupidrepo.RaylibGoGame) -v time
