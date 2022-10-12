package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestInfrastructureLayerDependencies(t *testing.T) {
	archtest.Package(t, emailImpl).ShouldNotDependOn(
		httpPresentation,
	)
	archtest.Package(t, loggerImpl).ShouldNotDependOn(
		httpPresentation,
	)
	archtest.Package(t, ratesImpl).ShouldNotDependOn(
		httpPresentation,
	)
	archtest.Package(t, reposImpl).ShouldNotDependOn(
		httpPresentation,
	)
}
