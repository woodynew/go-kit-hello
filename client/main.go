package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	client "github.com/woodynew/go-kit-hello/client/grpc"
	thrift1 "github.com/woodynew/go-kit-hello/client/thrift"
	service "github.com/woodynew/go-kit-hello/pkg/service"
	addthrift "github.com/woodynew/go-kit-hello/pkg/thrift/gen-go/addsvc"

	grpc1 "github.com/go-kit/kit/transport/grpc"

	"github.com/apache/thrift/lib/go/thrift"
	"google.golang.org/grpc"
)

func main() {
	thriftStart()

	// for i := 0; i < 10; i++ {
	// 	fmt.Println("start+" + cast.ToString(i))
	// 	// grpcStart()
	// 	thriftStart()
	// }
}
func grpcStart() {
	fmt.Printf("时间戳（毫秒）-start-grpcStart：%v;\n", time.Now().UnixNano()/1e6)

	fs := flag.NewFlagSet("hello", flag.ExitOnError)
	var grpcAddr = fs.String("grpc-addr", ":8082", "gRPC address of addsvc")
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	svc, err := client.New(conn, map[string][]grpc1.ClientOption{})
	if err != nil {
		panic(err)
	}

	r, err := svc.Foo(context.Background(), "hello")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", r)

	fmt.Printf("时间戳（毫秒）-end：%v;\n", time.Now().UnixNano()/1e6)

}

func thriftStart() {
	//控制台命令：go run main.go 1 2
	fs := flag.NewFlagSet("hello", flag.ExitOnError)
	var (
		thriftAddr     = fs.String("thrift-addr", ":8083", "Thrift address of addsvc")
		thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
		thriftBuffer   = fs.Int("thrift-buffer", 0, "0 for unbuffered")
		thriftFramed   = fs.Bool("thrift-framed", false, "true to enable framing")
		method         = fs.String("method", "sum", "sum, concat, foo")
	)

	fmt.Println("----------- start-thriftStart ----------------")
	fmt.Printf("时间戳（毫秒）-start：%v;\n", time.Now().UnixNano()/1e6)

	//os = 控制台输入参数
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	fs.Parse(os.Args[1:])

	fmt.Println(os.Args)
	fmt.Println(fs.Args())

	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	var (
		svc service.HelloService
		err error
	)

	var protocolFactory thrift.TProtocolFactory
	switch *thriftProtocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprintf(os.Stderr, "error: invalid protocol %q\n", *thriftProtocol)
		os.Exit(1)
	}
	var transportFactory thrift.TTransportFactory
	if *thriftBuffer > 0 {
		transportFactory = thrift.NewTBufferedTransportFactory(*thriftBuffer)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}
	if *thriftFramed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}
	transportSocket, err := thrift.NewTSocket(*thriftAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	transport, err := transportFactory.GetTransport(transportSocket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := transport.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer transport.Close()
	client := addthrift.NewAddServiceClientFactory(transport, protocolFactory)
	svc = thrift1.NewThriftClient(client)

	switch *method {
	case "sum":
		a, _ := strconv.ParseInt(fs.Args()[0], 10, 64)
		b, _ := strconv.ParseInt(fs.Args()[1], 10, 64)
		v, err := svc.Sum(context.Background(), int(a), int(b))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%d + %d = %d\n", a, b, v)

	case "concat":
		a := fs.Args()[0]
		b := fs.Args()[1]
		v, err := svc.Concat(context.Background(), a, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%q + %q = %q\n", a, b, v)

	case "foo":
		a := fs.Args()[0]
		v, err := svc.Foo(context.Background(), a)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%q , %q\n", a, v)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}

	fmt.Printf("时间戳（毫秒）-end：%v;\n", time.Now().UnixNano()/1e6)
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
