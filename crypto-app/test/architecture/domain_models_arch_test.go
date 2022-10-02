package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestThatDomainModelsLayerIsIndependent(t *testing.T) {
	archtest.Package(t, domainModels).ShouldNotDependOn(
		domainServices,
		application,
		httpPresentation,
		emailImpl,
		loggerImpl,
		ratesImpl,
		reposImpl,
	)
}
