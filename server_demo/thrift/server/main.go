package main

import (
	"context"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donscoco/gateway/server_demo/thrift/gen-go/thrift_gen"
	"log"
	"os"
)

const Addr = "127.0.0.1:30801"

type FormatDataImpl struct{}

func (fdi *FormatDataImpl) DoFormat(ctx context.Context, data *thrift_gen.Data) (r *thrift_gen.Data, err error) {
	var rData thrift_gen.Data
	rData.Text = Addr + " DoFormat"
	return &rData, nil
}

func main() {
	addr := flag.String("addr", Addr, "input addr")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	handler := &FormatDataImpl{}
	processor := thrift_gen.NewFormatDataProcessor(handler)
	serverSocket, err := thrift.NewTServerSocket(*addr)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverSocket, transportFactory, protocolFactory)
	fmt.Println("Running at:", *addr)
	if err := server.Serve(); err != nil {
		log.Fatalf(err.Error())
	}
}
