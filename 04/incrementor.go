// Пакет incrementor представляет собой потокобезопасный счетчик
// с возможностью установки предельного значения
// так же по заявлению разработчика этот идеальный пакет(насколько это в принципе возможно).

package incrementor

import (
	"math"
	"sync/atomic"
)

type Incrementor struct {
	i   uint32
	max uint32
}

// NewIncrementor возвращает экземпляр счетчика.
func NewIncrementor() *Incrementor {
	return &Incrementor{max: math.MaxUint32}
}

// GetNumber Возвращает текущее число. В самом начале это ноль.
func (i *Incrementor) GetNumber() uint32 {
	return atomic.LoadUint32(&i.i)
}

// GetMaximumValue Возвращает установленное максимальное значение.
func (i *Incrementor) GetMaximumValue() uint32 {
	return atomic.LoadUint32(&i.max)
}

// IncrementNumber Увеличивает текущее число на один. После каждого вызова этого
// метода GetNumber() будет возвращать число на один больше.
func (i *Incrementor) IncrementNumber() {
	if i.isMax() {
		i.setZero()
		return
	}

	atomic.AddUint32(&i.i, 1)
}

// SetMaximumValue Устанавливает максимальное значение текущего числа.
// Когда при вызове IncrementNumber() текущее число достигает
// этого значения, оно обнуляется, т.е. GetNumber() начинает
// снова возвращать ноль, и снова один после следующего
// вызова IncrementNumber() и так далее.
// По умолчанию максимум 4294967295.
// Если при смене максимального значения число начинает
// превышать максимальное значение, то число обнуляется.
func (i *Incrementor) SetMaximumValue(max uint32) {
	atomic.SwapUint32(&i.max, max)

	// NOTE: по условию задачи сразу же обнуляем счетчик (не очевидный момент)
	// при следующем IncrementNumber() мы получем 1 вместо 0
	if i.isMax() {
		i.setZero()
	}
}

// isMax проверяет достижение счетчиком максимального значения
func (i *Incrementor) isMax() bool {
	return atomic.LoadUint32(&i.i) >= atomic.LoadUint32(&i.max)
}

// setZero устанавливает нулевого значения
func (i *Incrementor) setZero() {
	atomic.SwapUint32(&i.i, 0)
}
