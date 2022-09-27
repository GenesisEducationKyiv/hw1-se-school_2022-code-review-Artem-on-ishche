package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
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
