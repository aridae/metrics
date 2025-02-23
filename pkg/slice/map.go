package slice

// MapBatch применяет функцию mapperFn к каждому элементу среза ins и собирает результаты в новый срез.
//
// Эта функция принимает срез элементов типа In и функцию mapperFn, которая преобразует каждый элемент в элемент типа Out.
// Если mapperFn возвращает ошибку, процесс прекращается и возвращается эта ошибка вместе с пустым срезом.
//
// Параметры:
//
//	ins - Входной срез элементов типа In.
//	mapperFn - Функция, принимающая элемент типа In и возвращающая элемент типа Out и возможную ошибку.
//
// Возвращаемое значение:
//
//	Срез результатов преобразования типа Out, если ошибки отсутствуют.
//	Ошибка, если mapperFn вернула ошибку для какого-либо элемента.
//
// Примеры:
//
//	MapBatch([]int{1, 2, 3}, func(x int) (int, error) { return x * 2, nil })  // Возвращает []int{2, 4, 6}, nil
//	MapBatch([]string{"apple", "banana", "cherry"}, strings.ToUpper)            // Возвращает []string{"APPLE", "BANANA", "CHERRY"}, nil
//
// Ошибки:
//
//	Функция возвращает ошибку, если mapperFn вернет ошибку для одного из элементов.
func MapBatch[In any, Out any](ins []In, mapperFn func(in In) (Out, error)) ([]Out, error) {
	outs := make([]Out, 0, len(ins))

	for _, in := range ins {
		out, err := mapperFn(in)
		if err != nil {
			return nil, err
		}

		outs = append(outs, out)
	}

	return outs, nil
}
