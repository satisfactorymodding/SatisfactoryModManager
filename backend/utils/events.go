package utils

import "slices"

type eventListener[D any] *func(D)

type EventDispatcher[D any] struct {
	listeners []eventListener[D]
}

func (ed *EventDispatcher[D]) On(f func(D)) func() {
	ed.listeners = append(ed.listeners, &f)
	return func() {
		ed.listeners = slices.DeleteFunc(ed.listeners, func(listener eventListener[D]) bool {
			return listener == &f
		})
	}
}

func (ed *EventDispatcher[D]) Once(f func(D)) {
	var after func()
	after = ed.On(func(data D) {
		f(data)
		after()
	})
}

func (ed *EventDispatcher[D]) Dispatch(data D) {
	for _, listener := range ed.listeners {
		if listener == nil {
			continue
		}
		(*listener)(data)
	}
	ed.listeners = slices.DeleteFunc(ed.listeners, func(listener eventListener[D]) bool {
		return listener == nil
	})
}
