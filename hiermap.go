/**
 * Copyright 2019 Innodev LLC. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package hiermap

import (
	"sync"
)

type HierMap struct {
	mux  *sync.Mutex
	maps []map[interface{}]interface{}
	size int
}

func New(size int) *HierMap {
	out := &HierMap{
		size: size,
		mux:  &sync.Mutex{},
		maps: make([]map[interface{}]interface{}, size),
	}

	for i := 0; i < size; i++ {
		out.maps[i] = map[interface{}]interface{}{}
	}

	return out
}

func (hm *HierMap) getMin(min []int) int {
	if len(min) == 0 {
		return 0
	}
	return min[0]
}

func (hm *HierMap) Delete(key interface{}, min ...int) {
	i := hm.getMin(min)
	hm.mux.Lock()
	defer hm.mux.Unlock()
	for ; i < hm.size; i++ {
		_, ok := hm.maps[i][key]
		if ok {
			delete(hm.maps[i], key)
			return
		}
	}
}

func (hm *HierMap) load(key interface{}, i int) (value interface{}, ok bool) {
	for ; i < hm.size; i++ {
		val, ok := hm.maps[i][key]
		if ok {
			return val, true
		}
	}
	return nil, false
}

func (hm *HierMap) store(key, value interface{}, i int) {
	hm.maps[i][key] = value
}

func (hm *HierMap) Load(key interface{}, min ...int) (value interface{}, ok bool) {
	hm.mux.Lock()
	defer hm.mux.Unlock()
	return hm.load(key, hm.getMin(min))
}

func (hm *HierMap) LoadOrStore(key, value interface{}, min ...int) (actual interface{}, loaded bool) {
	i := hm.getMin(min)
	hm.mux.Lock()
	defer hm.mux.Unlock()
	val, ok := hm.load(key, i)
	if ok {
		return val, true
	}
	hm.store(key, value, i)
	return value, false
}

func (hm *HierMap) Store(key, value interface{}, min ...int) {
	hm.mux.Lock()
	defer hm.mux.Unlock()
	hm.store(key, value, hm.getMin(min))
}
