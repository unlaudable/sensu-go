package types

import "errors"

// Validate returns an error if the CheckName and Subscription fields are not
// provided.
func (s *Silenced) Validate() error {
	checkNameErr := ValidateName(s.Check)
	subscriptionErr := ValidateName(s.Subscription)
	if checkNameErr != nil && subscriptionErr != nil {
		return errors.New("must provide check or subscription")
	}
	return nil
}

// GetOrg refers to the organization the check belongs to
func (s *Silenced) GetOrg() string {
	return s.Organization
}

// GetEnv refers to the organization the check belongs to
func (s *Silenced) GetEnv() string {
	return s.Environment
}

// FixtureSilenced returns a testing fixutre for a Silenced event struct.
func FixtureSilenced(checkName string) *Silenced {
	return &Silenced{
		Check: checkName,
	}
}