/**
 * Copyright 2019 Innodev LLC. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package hiermap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateValues(start, num int64) map[int64]int64 {
	out := make(map[int64]int64, num)
	for i := int64(0); i < num; i++ {
		out[i] = i + start
	}
	return out
}

func TestHiermap(t *testing.T) {
	limit := int64(1024)
	testData := []map[int64]int64{
		generateValues(0, limit),
		generateValues(limit*1, limit),
		generateValues(limit*2, limit),
		generateValues(limit*3, limit),
		generateValues(limit*4, limit),
	}

	testData2 := []map[int64]int64{
		generateValues(limit*10, limit*2),
		generateValues(limit*20, limit*2),
		generateValues(limit*30, limit*2),
		generateValues(limit*40, limit*2),
		generateValues(limit*50, limit*2),
	}

	hm := New(len(testData))

	for i := len(testData) - 1; i >= 0; i-- {
		for k, v := range testData[i] {
			hm.Store(k, v, i)
			v2, ok := hm.Load(k, i)
			assert.True(t, ok)
			assert.Equal(t, v, v2)
			if i > 0 {
				v2, ok := hm.Load(k, i-1)
				assert.True(t, ok)
				assert.Equal(t, v, v2)
			}
		}
	}
	awaitChan := make(chan bool)
	for i := range testData2 {
		go func(i int, dat map[int64]int64) {
			for k, v := range dat {
				val, ok := hm.LoadOrStore(k, v, i)
				if k >= limit {
					if ok {
						assert.True(t, val.(int64) > int64(limit*5))
					} else {
						assert.Equal(t, v, val)
					}
				} else {
					assert.True(t, ok)
					assert.NotEqual(t, v, val)
				}
			}
			awaitChan <- true
		}(i, testData2[i])
	}

	for range testData {
		<-awaitChan
	}

	for i := range testData {
		for k, v := range testData[i] {
			v2, ok := hm.Load(k, i)
			assert.True(t, ok)
			assert.Equal(t, v, v2)
		}
	}
}
