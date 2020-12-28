package zstdpool

import (
	"bytes"
	_ "log"
	"testing"

	"github.com/klauspost/compress/zstd"
)

var emptyZstdBlob = []byte{40, 181, 47, 253, 32, 0, 1, 0, 0}

func TestDecoders(t *testing.T) {
	t.Parallel()

	var pool DecoderPool

	numItems := 10
	items := make([]*zstd.Decoder, numItems*2, numItems*2)

	for i := 0; i < numItems; i++ {
		d, err := pool.Get(bytes.NewReader(emptyZstdBlob))
		if err != nil {
			t.Fatal(err)
		}
		if d == nil {
			t.Fatal("Expected non-nil *zstd.Decoder")
		}

		items[i] = d
	}

	if pool.available != 0 {
		t.Fatalf("expected pool to contain no items, found %d",
			pool.available)
	}
	for i := 0; i < numItems; i++ {
		pool.Put(items[i])
		items[i] = nil
	}

	if pool.available != numItems {
		t.Fatalf("expected pool to contain %d items, found %d",
			numItems, pool.available)
	}

	numItems *= 2

	for i := 0; i < numItems; i++ {
		d, err := pool.Get(bytes.NewReader(emptyZstdBlob))
		if err != nil {
			t.Fatal(err)
		}
		if d == nil {
			t.Fatal("Expected non-nil *zstd.Decoder")
		}

		items[i] = d
	}

	if pool.available != 0 {
		t.Fatalf("expected pool to contain no items, found %d",
			pool.available)
	}
	for i := 0; i < numItems; i++ {
		pool.Put(items[i])
		items[i] = nil
	}

	if pool.available != numItems {
		t.Fatalf("expected pool to contain %d items, found %d",
			numItems, pool.available)
	}
}

func TestEncoders(t *testing.T) {
	t.Parallel()

	var pool EncoderPool

	numItems := 10
	items := make([]*zstd.Encoder, numItems*2, numItems*2)

	for i := 0; i < numItems; i++ {
		d, err := pool.Get(nil)
		if err != nil {
			t.Fatal(err)
		}
		if d == nil {
			t.Fatal("Expected non-nil *zstd.Encoder")
		}

		items[i] = d
	}

	if pool.available != 0 {
		t.Fatalf("expected pool to contain no decoders, found %d",
			pool.available)
	}
	for i := 0; i < numItems; i++ {
		pool.Put(items[i])
		items[i] = nil
	}

	if pool.available != numItems {
		t.Fatalf("expected pool to contain %d items, found %d",
			numItems, pool.available)
	}

	numItems *= 2

	for i := 0; i < numItems; i++ {
		d, err := pool.Get(nil)
		if err != nil {
			t.Fatal(err)
		}
		if d == nil {
			t.Fatal("Expected non-nil *zstd.Encoder")
		}

		items[i] = d
	}

	if pool.available != 0 {
		t.Fatalf("expected pool to contain no items, found %d",
			pool.available)
	}
	for i := 0; i < numItems; i++ {
		pool.Put(items[i])
		items[i] = nil
	}

	if pool.available != numItems {
		t.Fatalf("expected pool to contain %d items, found %d",
			numItems, pool.available)
	}
}

func TestResizeDecoderPool(t *testing.T) {
	t.Parallel()

	var pool DecoderPool

	numItems := 10

	func() {
		for i := 0; i < numItems; i++ {
			d, err := pool.Get(nil)
			if err != nil {
				t.Fatal(err)
			}
			defer pool.Put(d)
		}
	}()

	if pool.available != numItems {
		t.Fatalf("expected %d items in the pool, found %d",
			numItems, pool.available)
	}

	resizeFunc := func(old int) int {
		return old - 2
	}

	expAvailable := numItems
	for i := 0; i < numItems/2; i++ {

		old, new, err := pool.Resize(resizeFunc)
		if err != nil {
			t.Fatal(err)
		}
		if old != expAvailable {
			t.Fatalf("expected initial pool size to be %d, found %d",
				expAvailable, old)
		}

		expAvailable -= 2

		if new != expAvailable {
			t.Fatalf("expected next pool size to be %d, found %d",
				expAvailable, new)
		}
	}
}

func TestResizeEncoderPool(t *testing.T) {
	t.Parallel()

	var pool EncoderPool

	numItems := 10

	func() {
		for i := 0; i < numItems; i++ {
			d, err := pool.Get(nil)
			if err != nil {
				t.Fatal(err)
			}
			defer pool.Put(d)
		}
	}()

	if pool.available != numItems {
		t.Fatalf("expected %d items in the pool, found %d",
			numItems, pool.available)
	}

	resizeFunc := func(old int) int {
		return old - 2
	}

	expAvailable := numItems
	for i := 0; i < numItems/2; i++ {

		old, new, err := pool.Resize(resizeFunc)
		if err != nil {
			t.Fatal(err)
		}
		if old != expAvailable {
			t.Fatalf("expected initial pool size to be %d, found %d",
				expAvailable, old)
		}

		expAvailable -= 2

		if new != expAvailable {
			t.Fatalf("expected next pool size to be %d, found %d",
				expAvailable, new)
		}
	}
}
