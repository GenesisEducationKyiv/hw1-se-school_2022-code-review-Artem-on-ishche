package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestApplicationLayerDependencies(t *testing.T) {
	archtest.Package(t, application).ShouldNotDependOn(
		httpPresentation,
		emailImpl,
		loggerImpl,
		ratesImpl,
		reposImpl,
	)
}
