package main

import (
	"context"
	"fmt"

	"github.com/iansmith/parigot-example/hello-world/g/greeting/v1"
	"github.com/iansmith/parigot/g/syscall/v1"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	lib "github.com/iansmith/parigot/lib/go"
)

func main() {
	ctx := pcontext.NewContextWithContainer(pcontext.GuestContext(context.Background()), "[hello-world]main")

	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "helloworld: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()

	myId := clientInit(ctx, []lib.MustRequireFunc{greeting.MustRequire})
	greetService := greeting.MustLocate(ctx, myId)

	req := &greeting.FetchGreetingRequest{
		Tongue: greeting.Tongue_French,
	}
	greetFuture := greetService.FetchGreeting(ctx, req)
	greetFuture.Method.Success(func(response *greeting.FetchGreetingResponse) {
		pcontext.Infof(ctx, "%s, world", response.Greeting)
		syscallguest.Exit(0)
	})
	greetFuture.Method.Failure(func(err greeting.GreetErr) {
		pcontext.Errorf(ctx, "failed to fetch greeting: %s", greeting.GreetErr_name[int32(err)])
		syscallguest.Exit(1)
	})
	clientOnlyRun(ctx)
}

func clientInit(ctx context.Context, requirement []lib.MustRequireFunc) id.ServiceId {

	myId := lib.MustRegisterClient(ctx)
	for _, f := range requirement {
		f(ctx, myId)
	}
	syscallguest.MustSatisfyWait(ctx, myId)

	launchreq := &syscall.LaunchRequest{
		ServiceId: myId.Marshal(),
	}
	_, err := syscallguest.Launch(launchreq)
	if err != syscall.KernelErr_NoError {
		panic(fmt.Sprintf("unable to launch client service: %s",
			syscall.KernelErr_name[int32(err)]))
	}

	return myId
}

var timeoutInMillis = int32(500)

func clientOnlyRun(ctx context.Context) {
	for {
		clientOnlyReadOneAndCall(ctx, nil, timeoutInMillis)
	}
}
func clientOnlyReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap,
	timeoutInMillis int32) syscall.KernelErr {
	req := syscall.ReadOneRequest{}

	// setup a request to read an incoming message
	req.TimeoutInMillis = timeoutInMillis
	req.HostId = lib.CurrentHostId().Marshal()
	resp, err := syscallguest.ReadOne(&req)
	if err != syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	lib.ExpireMethod(ctx)

	// is a promise being completed that was fulfilled somewhere else
	if r := resp.GetResolved(); r != nil {
		cid := id.UnmarshalCallId(r.GetCallId())
		lib.CompleteCall(ctx, cid, r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	return syscall.KernelErr_NoError
}
