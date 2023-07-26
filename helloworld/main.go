package main

import (
	"context"
	"runtime/debug"

	"github.com/iansmith/parigot-example/helloworld/g/greeting/v1"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

func main() {
	// Create a logging
	ctx := pcontext.NewContextWithContainer(pcontext.GuestContext(context.Background()), "[hello-world]main")

	// Logging system needs help with panics, so we trap and Dump() the log data.
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "helloworld: trapped a panic in the guest side: %v", r)
			debug.PrintStack()
		}
		pcontext.Dump(ctx)
	}()

	myId := lib.MustInitClient(ctx, []lib.MustRequireFunc{greeting.MustRequire})
	fut := lib.LaunchClient(ctx, myId)
	fut.Failure(func(err syscall.KernelErr) {
		pcontext.Errorf(ctx, "failed to launch the hello world service: %s", syscall.KernelErr_name[int32(err)])
		lib.ExitClient(ctx, 1, myId, "exiting due to launch failure", "unable to exit, forcing exit with os.Exit() after failure to Launch()")
	})

	fut.Success(func(resp *syscall.LaunchResponse) {
		pcontext.Infof(ctx, "hello world launched successfully")
		afterLaunch(ctx, resp, myId, fut)
	})

	// MustRunClient should never return.  Timeout in millis is used
	// for the question of how long should we "wait" for a network call
	// before doing something else.
	err := lib.MustRunClient(ctx, timeoutInMillis)

	// Should not happen.
	pcontext.Errorf(ctx, "failed inside run: %s", syscall.KernelErr_name[int32(err)])
	lib.ExitClient(ctx, 1, myId, "failed in MustRunClient", "failed trying to exit after failure in MustRunClient")
}

func afterLaunch(ctx context.Context, _ *syscall.LaunchResponse, myId id.ServiceId, fut *syscallguest.LaunchFuture) {
	greetService := greeting.MustLocate(ctx, myId)

	req := &greeting.FetchGreetingRequest{
		Tongue: greeting.Tongue_French,
	}

	// Make the call to the greeting service.
	greetFuture := greetService.FetchGreeting(ctx, req)

	// Handle positive outcome.
	greetFuture.Method.Success(func(response *greeting.FetchGreetingResponse) {
		pcontext.Infof(ctx, "%s, world", response.Greeting)
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, myId, "exiting after successful call to greeting.FetchGreeting",
			"failed trying to exit after success, so forcing exit with os.Exit()")
	})

	//Handle negative outcome.
	greetFuture.Method.Failure(func(err greeting.GreetErr) {
		pcontext.Errorf(ctx, "failed to fetch greeting: %s", greeting.GreetErr_name[int32(err)])
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, myId, "exiting because we failed to call greet.FetchGreeting",
			"tried to exit after failed call to greet.FetchGreeting, failed so forcing exit with os.Exit()")
	})

}

var timeoutInMillis = int32(500)
