package vcsview

// Represents commit author model
type Contributor struct {
	// Contributor name (if exists)
	name string

	// Contributor email (if exists)
	email string
}

// Get contributor name
func (c Contributor) Name() string {
	return c.name
}

// Get contributor email
func (c Contributor) Email() string {
	return c.email
}

func (c Contributor) String() string {
	if email := c.Email(); email != "" {
		return c.name + " <" + email + ">"
	}

	return c.name
}
