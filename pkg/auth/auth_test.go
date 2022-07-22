package auth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type authClientTestSuite struct {
	suite.Suite
}

func (s *authClientTestSuite) TestLogin() {}

func (s *authClientTestSuite) TestRevokeToken() {}

func (s *authClientTestSuite) TestRefreshToken() {}

func TestAuthClient(t *testing.T) {
	suite.Run(t, new(authClientTestSuite))
}
