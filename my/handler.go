package my

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"google.golang.org/grpc/codes"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {

	tr := global.TraceProvider().Tracer("Open-tracing-main_handler")

	client := http.DefaultClient
	ctx := distributedcontext.NewContext(context.Background())

	err := tr.WithSpan(ctx, "incoming call", // root span here
		func(ctx context.Context) error {

			// create child span
			ctx, childSpan := tr.Start(ctx, "backend call")
			childSpan.AddEvent(ctx, "making backend call")

			// create backend request
			req, _ := http.NewRequest("GET", "https://www.google.com", nil)

			// inject context
			ctx, req = httptrace.W3C(ctx, req)
			httptrace.Inject(ctx, req)

			// do request
			log.Printf("Sending request...\n")
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			_, err = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			// close child span
			childSpan.End()

			trace.SpanFromContext(ctx).SetStatus(codes.OK)
			log.Printf("got response: %v\n", res.Status)
			fmt.Printf("%v\n", "OK") //change to status code from backend
			return err
		})

	if err != nil {
		panic(err)
	}
}
