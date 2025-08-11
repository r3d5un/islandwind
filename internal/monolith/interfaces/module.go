package interfaces

import "context"

type Module interface {
	// Setup sets up the module using the context and resources from the
	// monolith. For example, initializing database models, message queues
	// and so on.
	//
	// Be aware that in case of injecting resources from another module,
	// no method call should be made within the Setup process. This can
	// cause segmentation faults as there is no guarantee that injected modules
	// have completed their own startup process.
	//
	// If resources are required from other modules as part of the startup
	// process, add a PostSetup or Startup method, and set it to run after
	// Setup has completed.
	Setup(ctx context.Context, mono Monolith)
	// Shutdown performs any necessary cleanup tasks before application
	// termination. Examples of such tasks could be closing channels,
	// connections created by the module and so on.
	//
	// Typically called with the defer keyword immediately after Setup and
	// PostSetup.
	Shutdown()
}
