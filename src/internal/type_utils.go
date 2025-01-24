package internal 

// s = s[0:] more efficient than setting to []string{}
// for repeated usage as it just reduces the slice 
// length without changing slice capacity

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}