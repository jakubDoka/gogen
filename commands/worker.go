//gogen_template
package commands

// I is template alias
type I = string

// O is template Alias
type O = string

// CAP is template constant
const CAP = 0

//<<< Worker<I, O, CAP, NM>

// NMWorker accepts data though channel and returns result from procedure
// also check for termination data with help of terminate argument
func NMWorker(rec chan O, procedure func(I) O, terminate func(I) bool, omit bool) chan I {
	c := make(chan I, CAP)
	closed := false
	go func() {
		done := []O{}
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
