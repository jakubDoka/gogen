package commands


// RWorker accepts data though channel and returns result from procedure
// also check for termination data with help of terminate argument
func RWorker(rec chan struct{}, procedure func(string) struct{}, terminate func(string) bool, omit bool) chan string {
	c := make(chan string, 30)
	closed := false
	go func() {
		done := []struct{}{}
		for {
			if !closed {
				select {
				case m := <-c:
					if terminate(m) {
						closed = true
						continue
					}
					done = append(done, procedure(m))
				default:
				}
			}

			if omit {
				if closed {
					return
				}
				continue
			}

			if len(done) != 0 {
				select {
				case rec <- done[0]:
					done = done[1:]
				default:
				}
			} else if closed {
				return
			}

		}
	}()
	return c
}


// TWorker accepts data though channel and returns result from procedure
// also check for termination data with help of terminate argument
func TWorker(rec chan []Template, procedure func(string) []Template, terminate func(string) bool, omit bool) chan string {
	c := make(chan string, 30)
	closed := false
	go func() {
		done := [][]Template{}
		for {
			if !closed {
				select {
				case m := <-c:
					if terminate(m) {
						closed = true
						continue
					}
					done = append(done, procedure(m))
				default:
				}
			}

			if omit {
				if closed {
					return
				}
				continue
			}

			if len(done) != 0 {
				select {
				case rec <- done[0]:
					done = done[1:]
				default:
				}
			} else if closed {
				return
			}

		}
	}()
	return c
}
