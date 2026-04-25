#!/usr/bin/env bash
set -euo pipefail

TOOLCHAIN_BIN="${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin"
ANDROID_API="${ANDROID_API:-21}"
OUT_DIR="${OUT_DIR:-android/libs/arm64-v8a}"
OUT_LIB="${OUT_DIR}/libgame.so"
CC_BIN="aarch64-linux-android${ANDROID_API}-clang"

if [[ ! -x "${TOOLCHAIN_BIN}/${CC_BIN}" ]]; then
  echo "error: missing compiler: ${TOOLCHAIN_BIN}/${CC_BIN}" >&2
  echo "hint: verify ANDROID_NDK_HOME points to the NDK root directory" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"

NATIVE_APP_GLUE_DIR="$(pwd)/raylib-go/raylib/external/android/native_app_glue"

PATH="${TOOLCHAIN_BIN}:${PATH}" \
CC="${CC_BIN}" \
CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
CGO_CFLAGS="-I${NATIVE_APP_GLUE_DIR} ${CGO_CFLAGS:-}" \
go build -C ./game/ -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libgame.so" \
-o "../${OUT_LIB}"

if [[ ! -s "${OUT_LIB}" ]]; then
  echo "error: build finished but output is missing or empty: ${OUT_LIB}" >&2
  exit 1
fi

echo "Built ${OUT_LIB}"

# installDebug quiet with minimal output
./gradlew installDebug --quiet
# then launch the app on the first connected device
adb shell am start -n "io.github.stupidrepo.RaylibGoGame/.GameNativeActivity"
