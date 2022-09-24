package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestPresentationLayerDependencies(t *testing.T) {
	archtest.Package(t, httpPresentation).ShouldNotDependOn(
		emailImpl,
		loggerImpl,
		ratesImpl,
		reposImpl,
	)
}
