/**
 * Tracks how many listeners there are at this moment.
 * @file ihr-listening.c
 */

#include <quickio.h>

#define PREFIX "/fuzzer"

struct _delayed {
	gint64 after;
	GDestroyNotify fn;
	void *data;
};

static GMutex _lock;
static GThread *_th = NULL;
static GSequence *_delayed = NULL;
static gboolean _is_running = TRUE;

static int _delayed_cmp(
	const void *a,
	const void *b,
	gpointer nothing G_GNUC_UNUSED)
{
	return ((struct _delayed*)a)->after - ((struct _delayed*)b)->after;
}

static void _delay(gint64 delay, GDestroyNotify fn, void *data)
{
	struct _delayed d = {
		.after = qev_monotonic + delay,
		.fn = fn,
		.data = data,
	};

	g_mutex_lock(&_lock);
	g_sequence_insert_sorted(_delayed, g_slice_copy(sizeof(d), &d),
							_delayed_cmp, NULL);
	g_mutex_unlock(&_lock);
}

static void _delayed_on_cb(void *info_)
{
	struct evs_on_info *info = info_;
	evs_on_cb(TRUE, info);
	evs_on_info_free(info);
}

static enum evs_status _delayed_on(const struct evs_on_info *info)
{
	_delay(QEV_SEC_TO_USEC(1), _delayed_on_cb, evs_on_info_copy(info, FALSE));
	return EVS_STATUS_HANDLED;
}

static void _purge_delayed()
{
	while (TRUE) {
		g_mutex_lock(&_lock);

		struct _delayed *d;
		GSequenceIter *iter = g_sequence_get_begin_iter(_delayed);
		if (g_sequence_iter_is_end(iter)) {
			g_mutex_unlock(&_lock);
			return;
		}

		d = g_sequence_get(iter);
		if (d->after < qev_monotonic) {
			g_mutex_unlock(&_lock);
			return;
		}

		g_sequence_remove(iter);
		g_mutex_unlock(&_lock);

		d->fn(d->data);
		g_slice_free1(sizeof(*d), d);
	}
}

static void* _run(void *nothing G_GNUC_UNUSED)
{
	while (_is_running) {
		_purge_delayed();
		g_usleep(QEV_MS_TO_USEC(100));
	}

	g_sequence_free(_delayed);
	_delayed = NULL;

	return NULL;
}

static gboolean _app_init()
{
	evs_add_handler(PREFIX, "/delayed", NULL, _delayed_on, NULL, TRUE);
	evs_add_handler(PREFIX, "/reject", NULL, evs_no_on, NULL, FALSE);

	_delayed = g_sequence_new(NULL);
	_th = g_thread_new("fuzzer_main", _run, NULL);

	return TRUE;
}

static gboolean _app_exit()
{
	_is_running = FALSE;

	g_thread_join(_th);
	_th = NULL;

	return TRUE;
}

QUICKIO_APP(
	_app_init,
	_app_exit);
