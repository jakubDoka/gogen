package parser


// SS is a string set
type SS map[string]bool

// Add adds value to set
func (n SS) Add(val string) {
n[val] = true
}

// Rem removes value from set
func (n SS) Rem(val string) (ok bool) {
ok = n[val]
delete(n, val)
return
}

// Join joins o with n
func (n SS) Join(o SS) {
for k := range o {
n[k] = true
}
}

