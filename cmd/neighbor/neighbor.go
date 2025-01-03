package neighbor

import (
	"bufio"
	"net"
	"os"

	cmdx "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/printer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts = struct {
	NamespaceWordlist string
	NamespaceList     []string
}{}

func init() {
	NeighborCmd.AddCommand(NeighborSvcCmd, NeighborPodCmd)
	cmdx.RootCmd.AddCommand(NeighborCmd)
	NeighborPodCmd.Flags().StringVar(&Opts.NamespaceWordlist, "ns-file", "", "namespace wordlist file")
	NeighborPodCmd.Flags().StringSliceVar(&Opts.NamespaceList, "ns", []string{}, "namespace list")
}

var NeighborCmd = &cobra.Command{
	Use:   "neighbor",
	Short: "neighbor is a tool to discover k8s neighbor pods",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var NeighborSvcCmd = &cobra.Command{
	Use:   "svc",
	Short: "neighbor is a tool to discover k8s neighbor pod with services",
	Run: func(cmd *cobra.Command, args []string) {
		ipNets, err := pkg.ParseStringToIPNet(cmdx.Opts.PodCidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		r := RunSvc(ipNets, cmdx.Opts.ThreadingNum)
		printer.PrintResult(r, cmdx.Opts.OutputFile)
	},
}

var NeighborPodCmd = &cobra.Command{
	Use:   "pod",
	Short: "neighbor is a tool to discover k8s pod and available ip in subnet (require k8s coredns with pod verified config)",
	Run: func(cmd *cobra.Command, args []string) {
		if !pkg.CheckPodVerified() {
			log.Fatalf("k8s coredns with pod verified config could not be set")
		}
		if Opts.NamespaceWordlist != "" {
			f, e := os.OpenFile(Opts.NamespaceWordlist, os.O_RDONLY, 0666)
			if e != nil {
				log.Fatalf("open file %v failed: %v", Opts.NamespaceWordlist, e)
			}
			defer f.Close()
			fileScanner := bufio.NewScanner(f)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				Opts.NamespaceList = append(Opts.NamespaceList, fileScanner.Text())
			}
		}
		log.Tracef("namespace list: %v", Opts.NamespaceList)
		ipNets, err := pkg.ParseStringToIPNet(cmdx.Opts.PodCidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		r := RunPod(Opts.NamespaceList, ipNets, cmdx.Opts.ThreadingNum)
		printer.PrintResult(r, cmdx.Opts.OutputFile)
	},
}

func RunPod(ns []string, net *net.IPNet, num int) (finalRecord define.Records) {
	scan := mutli.ScanNeighbor(ns, net, num)
	for r := range scan {
		finalRecord = append(finalRecord, r...)
	}
	if len(finalRecord) == 0 {
		log.Warn("ScanSubnet Found Nothing")
		return
	}
	return
}

func RunSvc(net *net.IPNet, num int) (finalRecord define.Records) {
	scan := mutli.ScanNeighborSvc(net, num)
	for r := range scan {
		finalRecord = append(finalRecord, r...)
	}
	if len(finalRecord) == 0 {
		log.Warn("ScanSubnet Found Nothing")
		return
	}
	return
}
