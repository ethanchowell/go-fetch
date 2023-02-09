package main

import (
	"github.com/ethanchowell/go-fetch/pkg/cmd"
	"k8s.io/klog/v2"
)

func main() {
	rootCmd := cmd.New()
	if err := rootCmd.Execute(); err != nil {
		klog.Fatalln(err)
	}
}
