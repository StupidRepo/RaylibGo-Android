#!/usr/bin/env bash
set -euo pipefail

# open logcat for package io.github.stupidrepo.android

adb logcat -d --pid=$(adb shell pidof -s io.github.stupidrepo.RaylibGoGame) -v time
