package io.github.stupidrepo.RaylibGoGame;

import android.app.NativeActivity;
import android.os.Build;
import android.os.Bundle;
import android.view.View;
import android.view.Window;
import android.view.WindowInsets;
import java.lang.annotation.Inherited;
import java.lang.annotation.Native;

public class GameNativeActivity extends NativeActivity {

	private volatile int insetLeft;
	private volatile int insetTop;
	private volatile int insetRight;
	private volatile int insetBottom;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		installInsetsListener();
		requestInsetsRefresh();
	}

	@Override
	public void onAttachedToWindow() {
		super.onAttachedToWindow();
		requestInsetsRefresh();
	}

	public int[] queryInsetsPx() {
		requestInsetsRefresh();
		return new int[] { insetLeft, insetTop, insetRight, insetBottom };
	}

	private void installInsetsListener() {
		final View decor = getDecorView();
		if (decor == null) {
			return;
		}

		decor.setOnApplyWindowInsetsListener(
			new View.OnApplyWindowInsetsListener() {
				@Override
				public WindowInsets onApplyWindowInsets(
					View v,
					WindowInsets insets
				) {
					updateInsets(insets);
					return v.onApplyWindowInsets(insets);
				}
			}
		);
	}

	private void requestInsetsRefresh() {
		final View decor = getDecorView();
		if (decor == null) {
			return;
		}

		decor.post(
			new Runnable() {
				@Override
				public void run() {
					if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
						updateInsets(decor.getRootWindowInsets());
					}
					decor.requestApplyInsets();
				}
			}
		);
	}

	private void updateInsets(WindowInsets insets) {
		if (insets == null) {
			return;
		}

		// Legacy insets API is available on API 21+ and returns px values.
		insetLeft = insets.getSystemWindowInsetLeft();
		insetTop = insets.getSystemWindowInsetTop();
		insetRight = insets.getSystemWindowInsetRight();
		insetBottom = insets.getSystemWindowInsetBottom();
	}

	private View getDecorView() {
		final Window window = getWindow();
		if (window == null) {
			return null;
		}
		return window.getDecorView();
	}
}
