package main

import (
	"context"
	"runtime/debug"

	"github.com/iansmith/parigot-example/hello-world/g/greeting/v1"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
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
		syscallguest.Exit(0)
	})

	//Handle negative outcome.
	greetFuture.Method.Failure(func(err greeting.GreetErr) {
		pcontext.Errorf(ctx, "failed to fetch greeting: %s", greeting.GreetErr_name[int32(err)])
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	})

	// MustRunClient should never return.  Timeout in millis is used
	// for the question of how long should we "wait" for a network call
	// before doing something else.
	err := lib.MustRunClient(ctx, timeoutInMillis)

	// Should not happen.
	pcontext.Fatalf(ctx, "failed inside run: %s", syscall.KernelErr_name[int32(err)])
}

var timeoutInMillis = int32(500)
