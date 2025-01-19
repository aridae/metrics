package slice

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
