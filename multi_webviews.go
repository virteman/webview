package webview

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"runtime"
	"time"
)

const (
	WEBVIEW_MASTER_ID = "WEBVIEW_MASTER"
)

type UIExec func(w WebView, args ...interface{}) bool

type Runner struct {
	Wv    WebView
	Uexec UIExec
}

type Mwebviews struct {
	webviews map[string]WebView
	toRun    chan *Runner
}

var sMWv *Mwebviews = &Mwebviews{
	webviews: make(map[string]WebView),
	toRun:    make(chan *Runner, 3),
}

func NewWindow(settings Settings) (mw *Mwebviews, err error) {
	_, err = sMWv.GetWebViewByID(settings.ID)
	if err != nil { //不存在
		if settings.ID == WEBVIEW_MASTER_ID {
			//创建主窗口
			sMWv.CreateMain(settings)
		} else {
			//客窗口
			_, err = sMWv.GetWebViewByID(WEBVIEW_MASTER_ID)
			if err != nil {
				err = errors.New("master window have to be created first")
				return
			}
			sMWv.Create(settings)
		}
		err = nil
	} else {
		err = errors.New("window ID already exists")
	}
	mw = sMWv
	return
}

// get web view by id
func (m *Mwebviews) GetWebViewByID(id string) (wv WebView, err error) {
	if w, ok := m.webviews[id]; ok {
		wv = w
	} else {
		err = errors.New("id not exist")
	}
	return
}

// run runner in the Main UI Thread
func (m *Mwebviews) RunInUIThread(id string, r UIExec) (err error) {
	wv, _ := m.GetWebViewByID(id)
	if r != nil {
		tRunner := &Runner{wv, r}
		m.toRun <- tRunner
	} else if err == nil {
		err = errors.New("UIExec cannot be nil")
	}
	return
}

// create window
func (m *Mwebviews) CreateWindow(settings Settings) (w WebView) {
	if runtime.GOOS == "linux" {
		m.RunInUIThread(settings.ID, func(w WebView, args ...interface{}) bool {
			w = m.Create(settings)
			return true
		})
	} else { //windows
		w = m.Create(settings)
		for w.Loop(true) {
		}
	}
	return
}

// create Main window
func (m *Mwebviews) CreateMain(settings Settings) (w WebView) {
	w = m.Create(settings)

	if runtime.GOOS == "linux" {
		//main loop when linux gtk
		for w.Loop(true) {
			select {
			case r := <-m.toRun:
				fmt.Printf("main thread id: %d\n", unix.Gettid())
				r.Uexec(r.Wv)
			case <-time.After(1 * time.Nanosecond):
			}
		}
	} else {
		//main loop when windows
		for w.Loop(true) {
		}
	}
	return
}

// create window
func (m *Mwebviews) Create(settings Settings) (w WebView) {
	webView := New(settings)
	m.webviews[settings.ID] = webView
	w = webView
	return
}
