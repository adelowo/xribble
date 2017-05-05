package xribble

type Option func(*XribbleDriver)

func BaseDir(dir string) Option {
	return func(x *XribbleDriver) {
		x.baseDir = dir
	}
}
