package acctelemetry

type accDataHolder[T AccGraphic | AccPhysics | AccStatic] struct {
	mmap *mmap
	data *T
}

func (d *accDataHolder[T]) Close() error {
	if d.mmap != nil {
		d.mmap.Close()
		d.mmap = nil
	}
	d.data = nil
	return nil
}
