package must

func NotFailF(fn func() error) {
	NotFail(fn())
}

func NotFail(err any) {
	if err != nil {
		panic(err)
	}
}
