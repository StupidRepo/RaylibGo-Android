//go:build android
// +build android

package platformspecifics

/*
#cgo LDFLAGS: -landroid -llog
#include <android/native_window.h>
#include <android_native_app_glue.h>
#include <jni.h>

extern struct android_app* GetAndroidApp();

static int ndk_window_width() {
    // we need a pointer to NativeWindow...
	// but we can get that with GetAndroidApp()!

	struct android_app* app = GetAndroidApp();
	if (!app || !app->window) return 0;

    return ANativeWindow_getWidth(app->window);
}
static int ndk_window_height() {
    struct android_app* app = GetAndroidApp();
    if (!app || !app->window) return 0;

    return ANativeWindow_getHeight(app->window);
}

static int jni_has_and_clear_exception(JNIEnv* env) {
    if (!env) return 1;
    if (!(*env)->ExceptionCheck(env)) return 0;
    (*env)->ExceptionClear(env);
    return 1;
}

static int jni_get_env(JavaVM* vm, JNIEnv** out_env, int* out_attached) {
    if (!vm || !out_env || !out_attached) return 0;

    *out_env = NULL;
    *out_attached = 0;

    jint env_status = (*vm)->GetEnv(vm, (void**)out_env, JNI_VERSION_1_6);
    if (env_status == JNI_EDETACHED) {
        if ((*vm)->AttachCurrentThread(vm, out_env, NULL) != JNI_OK || *out_env == NULL) {
            return 0;
        }
        *out_attached = 1;
        return 1;
    }

    return env_status == JNI_OK && *out_env != NULL;
}

static void jni_release_env(JavaVM* vm, int attached) {
    if (attached && vm) {
        (*vm)->DetachCurrentThread(vm);
    }
}

static jmethodID jni_get_activity_method(JNIEnv* env, jobject activity, jclass* out_cls, const char* name, const char* sig) {
    if (!env || !activity || !out_cls || !name || !sig) return NULL;

    *out_cls = (*env)->GetObjectClass(env, activity);
    if (!*out_cls || jni_has_and_clear_exception(env)) return NULL;

    jmethodID method = (*env)->GetMethodID(env, *out_cls, name, sig);
    if (!method || jni_has_and_clear_exception(env)) return NULL;

    return method;
}

static int jni_read_int_array4(JNIEnv* env, jintArray arr, int out_values[4]) {
    if (!env || !arr || !out_values) return 0;

    jsize len = (*env)->GetArrayLength(env, arr);
    if (jni_has_and_clear_exception(env) || len < 4) return 0;

    jint values[4] = {0, 0, 0, 0};
    (*env)->GetIntArrayRegion(env, arr, 0, 4, values);
    if (jni_has_and_clear_exception(env)) return 0;

    out_values[0] = values[0];
    out_values[1] = values[1];
    out_values[2] = values[2];
    out_values[3] = values[3];
    return 1;
}

static int ndk_get_insets(int out_insets[4]) {
    if (!out_insets) return 0;

    out_insets[0] = 0;
    out_insets[1] = 0;
    out_insets[2] = 0;
    out_insets[3] = 0;

    struct android_app* app = GetAndroidApp();
    if (!app || !app->activity || !app->activity->vm || !app->activity->clazz) return 0;

    JavaVM* vm = app->activity->vm;
    JNIEnv* env = NULL;
    int attached = 0;
    if (!jni_get_env(vm, &env, &attached)) return 0;

    int ok = 0;
    jobject activity = app->activity->clazz;
    jclass cls = NULL;
    jmethodID method = jni_get_activity_method(env, activity, &cls, "queryInsetsPx", "()[I");
    if (!method) goto cleanup_cls;

    jintArray arr = (jintArray)(*env)->CallObjectMethod(env, activity, method);
    if (arr && !jni_has_and_clear_exception(env)) {
        ok = jni_read_int_array4(env, arr, out_insets);
    } else {
        jni_has_and_clear_exception(env);
    }
    if (arr) (*env)->DeleteLocalRef(env, arr);

cleanup_cls:
    if (cls) (*env)->DeleteLocalRef(env, cls);

done:
    jni_release_env(vm, attached);
    return ok;
}

#include <android/log.h>
#include <stdlib.h>

static int ndk_log_write(int prio, const char* tag, const char* text) {
	return __android_log_write(prio, tag, text);
}
*/
import "C"
import (
	"unsafe"
)

type Platform struct{}

const (
	GLSLVersion = 100
)

func (p *Platform) GetOS() PlatformEnum {
	return PlatformAndroid
}

func (p *Platform) GetWindowSize() (int32, int32) {
	width, height := int32(C.ndk_window_width()), int32(C.ndk_window_height())
	return width, height
}
func (p *Platform) GetInsets() Insets {
	var insets [4]C.int
	if C.ndk_get_insets(&insets[0]) == 0 {
		p.LogIt(AndroidLogWarn, NDKTag, "Insets unavailable via JNI; using zeros")
		return Insets{0, 0, 0, 0}
	}

	return Insets{
		Left:   int32(insets[0]),
		Top:    int32(insets[1]),
		Right:  int32(insets[2]),
		Bottom: int32(insets[3]),
	}
}

func (a *Platform) LogIt(priority int, tag string, text string) int {
	cTag := C.CString(tag)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cText))

	return int(C.ndk_log_write(C.int(priority), cTag, cText))
}
