package fxecumenos

import "github.com/blang/semver/v4"

type ServiceName string

type Version semver.Version

func NewVersion(s string) Version {
	return Version(semver.MustParse(s))
}
