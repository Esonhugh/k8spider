package admission

import (
	"errors"
	"fmt"
	"os"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/pkg/admission-webhook/reviewer"
	"github.com/esonhugh/k8spider/pkg/admission-webhook/reviewer/defaultResource"
	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/admission/v1"
)

var Opts = struct {
	As                 string
	Group              []string
	Indent             int
	Action             string
	SendToController   bool
	ControllerEndpoint string

	FileContent [][]byte

	UseDefaultContent string
}{}

func init() {
	AdmitCmd.Flags().StringVarP(&Opts.As, "as", "a", "", "as username ")
	AdmitCmd.Flags().StringSliceVarP(&Opts.Group, "group", "g", []string{}, "as group")
	AdmitCmd.Flags().IntVarP(&Opts.Indent, "indent", "I", 2, "indent")
	AdmitCmd.Flags().StringVarP(&Opts.Action, "action", "A", "create", "action")
	AdmitCmd.Flags().BoolVarP(&Opts.SendToController, "send-to-controller", "s", false, "send to controller")
	AdmitCmd.Flags().StringVarP(&Opts.ControllerEndpoint, "controller-endpoint", "e", "", "controller endpoint")
	AdmitCmd.Flags().StringVarP(&Opts.UseDefaultContent, "use-default-content", "C", "", "use default resource content if not specified")
	command.RootCmd.AddCommand(AdmitCmd)
}

var AdmitCmd = &cobra.Command{
	Use:   "admit [file]",
	Short: "admit is a tool to testing kubernetes admission controllers",
	Long:  "admit is a tool to testing kubernetes admission controllers",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if Opts.SendToController {
			if Opts.ControllerEndpoint == "" {
				return errors.New("if send to controller is ture, controller-endpoint can't be empty")
			}
		}
		if Opts.UseDefaultContent != "" {
			switch Opts.UseDefaultContent {
			case "pod":
				Opts.FileContent = append(Opts.FileContent, []byte(defaultResource.Resources.Pod))
			case "deployment", "deploy":
				Opts.FileContent = append(Opts.FileContent, []byte(defaultResource.Resources.Deployment))
			case "ingress", "ing":
				Opts.FileContent = append(Opts.FileContent, []byte(defaultResource.Resources.Ingress))
			case "ingress-tls", "ing-tls":
				Opts.FileContent = append(Opts.FileContent, []byte(defaultResource.Resources.TLSIngress))
			default:
				return errors.New("use-default-content only support pod, deployment(deploy), ingress(ing), ingress-tls(ing-tls)")
			}
			return nil
		}
		if len(args) == 0 {
			return errors.New("file args can't be empty")
		}
		for _, file := range args {
			f, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			Opts.FileContent = append(Opts.FileContent, f)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, input := range Opts.FileContent {
			content, err := reviewer.CreateAdmissionReviewRequest(input, Opts.Action, Opts.As, Opts.Group, Opts.Indent)
			if err != nil {
				log.Errorf("create admission review request failed: %v", err)
				return
			}
			log.Debugf("admission review request: %v", string(content))
			if Opts.SendToController {
				log.Tracef("sending to the controller: %v", Opts.ControllerEndpoint)
				gp := gout.POST(Opts.ControllerEndpoint).SetJSON(content)
				if command.Opts.Verbose > 0 {
					gp.Debug(true)
				}
				var resp v1.AdmissionReview
				err := gp.BindJSON(&resp).Do()
				if err != nil {
					log.Errorf("send to controller failed: %v", err)
				}
				if resp.Response == nil {
					log.Errorf("get empty admission reviewed response")
				}
				log.Infof("get admission reviewed response")
				log.Infof("Allowed: %v", resp.Response.Allowed)
				log.Infof("UID: %v", resp.Response.UID)
				log.Infof("PatchType: %v", *resp.Response.PatchType)
				log.Infof("Patch: %v", string(resp.Response.Patch))
				for k, v := range resp.Response.AuditAnnotations {
					log.Infof("AuditAnnotation[%v]: %v", k, v)
				}
				for _, w := range resp.Response.Warnings {
					log.Infof("Warning: %v", w)
				}
				return
			} else {
				fmt.Println(string(content))
				if command.Opts.OutputFile != "" {
					err := os.WriteFile(command.Opts.OutputFile, content, 0644)
					if err != nil {
						log.Errorf("write to file failed: %v", err)
					}
				}
			}
		}
	},
}
