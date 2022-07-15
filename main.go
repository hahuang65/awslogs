package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"git.sr.ht/~hwrd/awslogs/internal/aws"
	"git.sr.ht/~hwrd/awslogs/internal/tui"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func PrintDescription() {
	commandName := path.Base(os.Args[0])

	fmt.Fprintf(flag.CommandLine.Output(),
		`%s is client for viewing AWS logs

`, commandName)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", commandName)
}

func main() {
	w := os.Stdout
	if err := run(os.Args, w); err != nil {
		fmt.Fprintf(w, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, w io.Writer) error {
	var (
		showHelp bool
	)

	flag.Usage = func() {
		PrintDescription()
		flag.PrintDefaults()
	}

	flag.BoolVar(&showHelp, "h", false, "shows this help guide")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	cloudWatchLogsClient := cloudwatchlogs.NewFromConfig(cfg)
	s := aws.Service{CloudWatchLogsClient: cloudWatchLogsClient}

	return tui.Start(s)
}
