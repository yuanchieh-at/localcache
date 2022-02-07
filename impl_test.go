package localcache

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type localCacheSuite struct {
	suite.Suite
	c *cache
}

func (s *localCacheSuite) SetupTest() {
	s.c = New()
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(localCacheSuite))
}

func (s *localCacheSuite) TestGet() {
	tests := []struct{
		Desc string
		SetupTest func(string)
		Key string
		ExpError error
		ExpResult interface{}
	} {
		{
			Desc: "not existed",
			Key: "not existed",
			ExpError: NewKeyNotFound("not existed"),
			ExpResult: nil,
		},
		{
			Desc: "get result",
			SetupTest: func(desc string) {
				_ = s.c.Set("get result", "value")
			},
			Key: "get result",
			ExpError: nil,
			ExpResult: "value",
		},
		{
			Desc: "key expired",
			Key: "expired",
			SetupTest: func(desc string) {
				timeNow = func() time.Time {
					return time.Now().Add(ttl * -1)
				}
				_ = s.c.Set("expired", "expired")
				timeNow = func() time.Time {
					return time.Now()
				}
			},
			ExpError: NewKeyNotFound("expired"),
			ExpResult: nil,
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.SetupTest(t.Desc)
		}

		value, err := s.c.Get(t.Key)
		s.Require().Equal(t.ExpError, err, t.Desc)
		if err == nil {
			s.Require().Equal(t.ExpResult, value, t.Desc)
		}
	}
}

func (s *localCacheSuite) TestSet() {
	tests := []struct{
		Desc string
		SetupTest func(string)
		Key string
		ExpResult interface{}
	}{
		{
			Desc:      "set value if key is not existed",
			Key:       "not existed",
			SetupTest: func(desc string) {
				_ = s.c.Set("not existed", "value")
			},
			ExpResult: "value",
		},
		{
			Desc: "override if key exits",
			SetupTest: func(desc string) {
				_ = s.c.Set("override", "value")
				_ = s.c.Set("override", "new value")
			},
			Key:       "override",
			ExpResult: "new value",
		},
	}

	for _, t := range tests {
		if t.SetupTest != nil {
			t.SetupTest(t.Desc)
		}

		value, err := s.c.Get(t.Key)
		if err == nil {
			s.Require().Equal(t.ExpResult, value, t.Desc)
		}
	}
}