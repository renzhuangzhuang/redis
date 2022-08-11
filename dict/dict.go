package dict

import (
	"errors"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/dchest/siphash"
)

// redis 字典实现
// 这边参考一下网上 对于key value 不做限制 使用interface{}来实现

const (
	_initialHashtableSize uint64 = 4
)

// 字典结构
type Dict struct {
	hashTables []*hashTable
	rehashIdx  int64
	iterators  uint64 //迭代器
}

//哈希表
type hashTable struct {
	dictEntry []*entry //哈希表数组
	size      uint64   // 哈希表大小
	sizemask  uint64   //哈希表掩码 计算索引值
	used      uint64   // 表示已有节点的数量
}

//哈希表节点
type entry struct {
	key, value interface{} // 表示我可以传入任意的数据类型
	next       *entry      // 指向下一个哈希表节点
}

// 初始化一个dict
func NewDict() *Dict {
	return &Dict{
		// 根据书上介绍
		// dict 准备两个hashtable 默认用第一个 第二用于rehash操作
		// rehashIdx 默认为 -1 不行 rehash操作
		hashTables: []*hashTable{{}, {}},
		rehashIdx:  -1,
		iterators:  0,
	}
}

// 计算字典中元素个数
func (d *Dict) Len() uint64 {
	var _len uint64
	for _, ht := range d.hashTables {
		_len += ht.used
	}
	return _len
}

// 返回字典容量
func (d *Dict) Cap() uint64 {
	if d.isRehashing() {
		return d.hashTables[1].size
	}
	return d.hashTables[0].size
}

// SipHash 算法
func SipHash(key interface{}) uint64 {
	// 第一步为了能够使用import出来的函数
	// 将接口类型数据进行转换为[]byte
	//keybyte := key.([]byte)
	// 采用分支判断方法
	var data []byte
	switch iv := key.(type) {
	case string:
		data = []byte(iv)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		data = []byte(fmt.Sprintf("%d", iv))
	default:
		panic(fmt.Sprintf("key type '%s' is not supported", reflect.TypeOf(key).String()))
	}
	value := siphash.Hash(uint64(12), uint64(22), data)
	return value

}

// 判断是否需要进行渐进性哈希
func (d *Dict) isRehashing() bool {
	if d.rehashIdx < 0 {
		return false
	} else {
		return true
	}
}

func (d *Dict) rehashStep() {
	if d.iterators == 0 {
		d.rehash(1)
	}

}

// 哈希表查找指定键
func (d *Dict) keyIndex(key interface{}) (idx uint64, existed *entry) {
	hash := SipHash(key) // 得到一个uint64的数字
	for i := 0; i < 2; i++ {
		ht := d.hashTables[i]
		idx = ht.sizemask & hash
		for ent := ht.dictEntry[idx]; ent != nil; ent = ent.next {
			if ent.key == key {
				return idx, ent
			}
		}
		if !d.isRehashing() {
			break

		}
	}
	return idx, nil
}

//查找键值对
func (d *Dict) Load(key interface{}) (value interface{}, ok bool) {
	if d.isRehashing() {
		d.rehashStep()
	}

	_, existed := d.keyIndex(key)
	if existed != nil {
		return existed.value, true
	}

	return nil, false
}

// 存储键值对
func (d *Dict) Store(key interface{}, value interface{}) {
	ent, loaded := d.loadOrStore(key, value)
	if loaded {
		ent.value = value
	}
}

// LoadOrStore 如果 key 存在于字典中，则直接返回其对应的值。
// 否则，该函数会将给定的值添加的字典中，并将给定的默认值返回。
// 如果能够在字典中成功查找的给定的 key，则 loaded 返回 true，
// 否则返回 false。
func (d *Dict) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	ent, loaded := d.loadOrStore(key, value)
	if loaded {
		return ent.value, true
	} else {
		return value, false
	}
}
func (d *Dict) loadOrStore(key, value interface{}) (ent *entry, loaded bool) {
	if d.isRehashing() {
		d.rehashStep()
	}

	_ = d.expandIfNeeded() //扩容
	idx, existed := d.keyIndex(key)
	ht := d.hashTables[0]

	if d.isRehashing() {
		ht = d.hashTables[1]
	}

	if existed != nil {
		return existed, true
	} else {
		// 不存在key时候 需要在dictEntry中添加新的entry
		// 键冲突情况采用链地址法，采用头插法
		entry := &entry{
			key:   key,
			value: value,
			next:  ht.dictEntry[idx],
		}
		ht.dictEntry[idx] = entry
		ht.used++
	}
	return nil, false
}

// 删除键值对
func (d *Dict) Delete(key interface{}) {
	if d.Len() == 0 {
		return
	}

	if d.isRehashing() {
		d.rehashStep()
	}

	hash := SipHash(key)
	for i := 0; i < 2; i++ {
		ht := d.hashTables[i]
		idx := ht.sizemask & hash
		var prevEntry *entry
		for ent := ht.dictEntry[idx]; ent != nil; ent = ent.next {
			if ent.key == key {
				if prevEntry != nil {
					prevEntry.next = ent.next
				} else {
					ht.dictEntry[idx] = ent.next
				}

				ent.next = nil
				ht.used--
				return
			}
			prevEntry = ent
		}
		if !d.isRehashing() {
			break
		}
	}

}

func (d *Dict) expandIfNeeded() error {
	if d.isRehashing() {
		//正在进行rehash
		return nil
	}

	if d.hashTables[0].size == 0 {
		// 第一次扩容， 需要一定的空间
		return d.resizeTo(_initialHashtableSize)
	}

	// 根据负载因子判断是否需要进行扩容

	if d.hashTables[0].used == d.hashTables[0].size {
		return d.resizeTo(d.hashTables[0].used * 2)
	}

	return nil
}

func (d *Dict) resizeTo(size uint64) error {
	if d.isRehashing() || d.hashTables[0].used > size {
		return errors.New("failed to resize")
	}

	size = d.nextPower(size)
	if size == d.hashTables[0].size {
		return nil
	}

	var ht *hashTable
	if d.hashTables[0].size == 0 {
		// 第一次扩容
		ht = d.hashTables[0]
	} else {
		ht = d.hashTables[1]
		// 开始进一步扩容
		d.rehashIdx = 0
	}

	ht.size = size
	ht.sizemask = size - 1
	ht.dictEntry = make([]*entry, ht.size)
	return nil
}

func (d *Dict) nextPower(size uint64) uint64 {
	/* if size >= math.MaxUint64 {
		return math.MaxUint64
	} */

	i := _initialHashtableSize
	for i < size {
		i <<= 1 // i*= 2
	}

	return i
}

// Resize 让字典扩容或者缩容一定大小
func (d *Dict) Resize() error {
	if d.isRehashing() {
		return errors.New("dict is rehashing")
	}

	size := d.hashTables[0].used
	if size < _initialHashtableSize {
		size = _initialHashtableSize
	}

	return d.resizeTo(size)
}

//渐进式rehash
func (d *Dict) rehash(steps uint64) (finished bool) {
	if !d.isRehashing() {
		return true
	}

	maxDictEntryMeets := 10 * steps
	src, dst := d.hashTables[0], d.hashTables[1]
	for ; steps > 0 && src.used != 0; steps-- {
		for src.dictEntry[d.rehashIdx] == nil {
			d.rehashIdx++
			maxDictEntryMeets--
			if maxDictEntryMeets <= 0 {
				return false
			}
		}
		for ent := src.dictEntry[d.rehashIdx]; ent != nil; {
			next := ent.next
			idx := SipHash(ent.key) & dst.sizemask
			ent.next = dst.dictEntry[idx]
			dst.dictEntry[idx] = ent
			src.used--
			dst.used++
			ent = next
		}
		src.dictEntry[d.rehashIdx] = nil
		d.rehashIdx++
	}

	if src.used == 0 {
		d.hashTables[0] = dst
		d.hashTables[1] = &hashTable{}
		d.rehashIdx = -1
		return true
	}

	return false
}

// 实现dict中迭代器功能-- 借鉴 gitee上作者 链接：https://gitee.com/ifaceless/go-redis-dict/blob/master/iter.go
// iterator 实现了一个对字典的迭代器。
// 不过考虑到我们将为字典提供 `Range` 方法，故该迭代器就不往外暴露了。
type iterator struct {
	d                  *Dict
	tableIndex         int
	safe               bool
	fingerprint        int64
	entry              *entry
	bucketIndex        uint64
	waitFirstIteration bool
}

func newIterator(d *Dict, safe bool) *iterator {
	return &iterator{
		d:                  d,
		safe:               safe,
		waitFirstIteration: true,
	}
}

// next 会依次扫描字典中哈希表的所有 buckets，并将其中的 entry 一一返回。
// 如果字典正在 rehash，那么会在扫描完哈希表 1 后，继续扫描哈希表 2。需要
// 注意的是，如果在迭代期间，继续向字典中添加数据可能没法被扫描到！
func (it *iterator) next() *entry {
	for {
		if it.entry == nil {
			if it.waitFirstIteration {
				// 第一次迭代，要做点特别的事情~
				if it.safe {
					// 告诉 dict，有正在运行的安全迭代器，进而阻止某些操作时的 rehash 操作
					it.d.iterators++
				} else {
					it.fingerprint = it.d.fingerprint()
				}
				it.waitFirstIteration = false
			}

			ht := it.d.hashTables[it.tableIndex]
			if it.bucketIndex >= ht.size {
				if !it.d.isRehashing() || it.tableIndex != 0 {
					return nil
				}

				// 切换到第二个哈希表继续扫描
				it.tableIndex = 1
				it.bucketIndex = 0
				ht = it.d.hashTables[1]
			}

			it.entry = ht.dictEntry[it.bucketIndex]
			it.bucketIndex++
		} else {
			it.entry = it.entry.next
		}

		if it.entry != nil {
			return it.entry
		}
	}
}

func (it *iterator) release() {
	if it.safe {
		it.d.iterators--
	} else {
		fp := it.d.fingerprint()
		if fp != it.fingerprint {
			panic("operations like 'LoadOrStore', 'Load' or 'Delete' are not safe for an unsafe iterator")
		}
	}
}

func (d *Dict) rangeDict(fn func(key, value interface{}) bool, safe bool) {
	it := newIterator(d, safe)
	defer it.release()

	for {
		if ent := it.next(); ent != nil {
			if !fn(ent.key, ent.value) {
				break
			}
		} else {
			break
		}
	}
}

// fingerprint 计算出字典指纹，相当于给字典某刻状态盖个戳。
// 在非安全模式下使用迭代器时，原则上是不允许执行查找、插入和删除
// 操作的，否则可能引起 rehash，导致迭代结果可能会重复扫描到某
// 些 keys；所以这里会在迭代结束后再次计算下字典状态，保证前后
// 指纹相等，否则需要告知用户进行了不当的操作~
// 算法参考 redis/src/dict.c#dictFingerprint
func (d *Dict) fingerprint() int64 {
	metas := []int64{
		// meta of table 0
		int64(uintptr(unsafe.Pointer(&d.hashTables[0].dictEntry))),
		int64(d.hashTables[0].size),
		int64(d.hashTables[0].used),
		// meta of table 1
		int64(uintptr(unsafe.Pointer(&d.hashTables[1].dictEntry))),
		int64(d.hashTables[1].size),
		int64(d.hashTables[1].used),
	}

	var hash int64
	for _, meta := range metas {
		hash += meta
		// 使用 Tomas Wang 64 位整数 hash 算法
		hash = (hash << 21) - hash - 1
		hash = hash ^ (hash >> 24)
		hash = (hash + (hash << 3)) + (hash << 8) // hash * 256
		hash = hash ^ (hash >> 14)
		hash = (hash + (hash << 2)) + (hash << 4) // hash * 21
		hash = hash ^ (hash >> 28)
		hash = hash + (hash << 31)
	}

	return hash
}

// 下面实现dict中的range功能
// Range 以非安全的方式进行迭代，意味着在迭代期间不允许对字典执行额外的
// 操作，以免引起 rehash，导致重复扫描一些键。
// 用户传入的 `fn` 回调，可以通过返回 false 指示迭代器停止工作。
// 迭代完毕后，内部迭代器在释放时会自动对比当前的字典指纹和迭代器指纹是否
// 一致，如果不一致，表明执行了禁止的操作，进而引起 panic。
func (d *Dict) Range(fn func(key, value interface{}) bool) {
	d.rangeDict(fn, false)
}

// Range 以安全的方式进行迭代，可以在迭代期间执行 Load, Store 等操作，
// 它会在执行这些操作时，阻止字典进行 rehash 操作。但不保证新加入的键值
// 一定能够被扫描到。
// 用户传入的 `fn` 回调，可以通过返回 false 指示迭代器停止遍历。
func (d *Dict) RangeSafely(fn func(key, value interface{}) bool) {
	d.rangeDict(fn, true)
}

// RehashForAWhile 执行 rehash 一段时间。
func (d *Dict) RehashForAWhile(duration time.Duration) int64 {
	tm := time.NewTimer(duration)
	defer tm.Stop()

	var rehashes int64
	for {
		select {
		case <-tm.C:
			return rehashes
		default:
			if d.rehash(100) {
				return rehashes
			}
			rehashes += 100
		}
	}
}
