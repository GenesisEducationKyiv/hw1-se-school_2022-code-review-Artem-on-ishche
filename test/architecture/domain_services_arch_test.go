package architecture

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestDomainServicesLayerDependencies(t *testing.T) {
	archtest.Package(t, domainServices).ShouldNotDependOn(
		application,
		httpPresentation,
		emailImpl,
		loggerImpl,
		ratesImpl,
		reposImpl,
	)
}
