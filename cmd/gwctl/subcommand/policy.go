package subcommand

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cmdutil "github.ibm.com/mbg-agent/cmd/util"
	"github.ibm.com/mbg-agent/pkg/admin"
	event "github.ibm.com/mbg-agent/pkg/controlplane/eventManager"
	"github.ibm.com/mbg-agent/pkg/policyEngine"
)

const (
	acl = "acl"
	lb  = "lb"
)

// PolicyCreateOptions is the command line options for 'create policy'
type policyCreateOptions struct {
	myID       string
	pType      string
	serviceSrc string
	serviceDst string
	gwDest     string
	priority   int
	action     int
	policy     string
}

// PolicyCreateCmd - create a new policy - TODO update this command after integration.
func PolicyCreateCmd() *cobra.Command {
	o := policyCreateOptions{}
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "An applyPolicy command send the gateway the policy for dedicated service",
		Long:  `An applyPolicy command send the gateway the policy for dedicated service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run()
		},
	}
	o.addFlags(cmd.Flags())
	cmdutil.MarkFlagsRequired(cmd, []string{"type"})

	return cmd
}

// addFlags registers flags for the CLI.
func (o *policyCreateOptions) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.myID, "myid", "", "gwctl ID")
	fs.StringVar(&o.pType, "type", "", "Policy agent command (For now, acl,lb)")
	fs.StringVar(&o.serviceSrc, "serviceSrc", "*", "Name of Source Service (* for wildcard)")
	fs.StringVar(&o.serviceDst, "serviceDst", "*", "Name of Dest Service (* for wildcard)")
	fs.StringVar(&o.gwDest, "gwDest", "*", "Name of gateway the dest service belongs to (* for wildcard)")
	fs.IntVar(&o.priority, "priority", 0, "Priority of the acl rule (0 -> highest)")
	fs.IntVar(&o.action, "action", 0, "acl 0 -> allow, 1 -> deny")
	fs.StringVar(&o.policy, "policy", "random", "lb policy: random, ecmp, static")
}

// run performs the execution of the 'create policy' subcommand
func (o *policyCreateOptions) run() error {
	g, err := admin.GetClientFromID(o.myID)
	if err != nil {
		return err
	}
	switch o.pType {
	case acl:
		return g.SendACLPolicy(o.serviceSrc, o.serviceDst, o.gwDest, o.priority, event.Action(o.action), admin.Add)
	case lb:
		return g.SendLBPolicy(o.serviceSrc, o.serviceDst, policyEngine.PolicyLoadBalancer(o.policy), o.gwDest, admin.Add)

	default:
		return fmt.Errorf("Unknown policy type")
	}
}

// PolicyDeleteOptions is the command line options for 'delete policy'
type policyDeleteOptions struct {
	myID       string
	pType      string
	serviceSrc string
	serviceDst string
	gwDest     string
	policy     string
}

// PolicyDeleteCmd - delete a policy command - TODO change after the policy integration.
func PolicyDeleteCmd() *cobra.Command {
	o := policyDeleteOptions{}
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Delete service policy from gateway.",
		Long:  `Delete service policy from gateway.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run()
		},
	}
	o.addFlags(cmd.Flags())
	cmdutil.MarkFlagsRequired(cmd, []string{"type"})

	return cmd
}

// addFlags registers flags for the CLI.I.
func (o *policyDeleteOptions) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.myID, "myid", "", "gwctl ID")
	fs.StringVar(&o.pType, "type", "", "Policy agent command (For now, acl,lb)")
	fs.StringVar(&o.serviceSrc, "serviceSrc", "*", "Name of Source Service (* for wildcard)")
	fs.StringVar(&o.serviceDst, "serviceDst", "*", "Name of Dest Service (* for wildcard)")
	fs.StringVar(&o.gwDest, "gwDest", "*", "Name of gateway the dest service belongs to (* for wildcard)")
	fs.StringVar(&o.policy, "policy", "random", "lb policy: random, ecmp, static")
}

// run performs the execution of the 'delete policy' subcommand
func (o *policyDeleteOptions) run() error {
	priority := 0 //Doesn't matter when deleting a rule
	action := 0   //Doesn't matter when deleting a rule
	g, err := admin.GetClientFromID(o.myID)
	if err != nil {
		return err
	}
	switch o.pType {
	case acl:
		err = g.SendACLPolicy(o.serviceSrc, o.serviceDst, o.gwDest, priority, event.Action(action), admin.Del)
	case lb:
		err = g.SendLBPolicy(o.serviceSrc, o.serviceDst, policyEngine.PolicyLoadBalancer(o.policy), o.gwDest, admin.Del)
	default:
		return fmt.Errorf("Unknown policy type")
	}
	return err
}

// PolicyGetOptions is the command line options for 'get policy'
type policyGetOptions struct {
	myID string
}

// PolicyGetCmd - get a policy command
func PolicyGetCmd() *cobra.Command {
	o := policyGetOptions{}
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Get policy list from the GW",
		Long:  `Get policy list from the GW (ACL and LB)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run()
		},
	}
	o.addFlags(cmd.Flags())

	return cmd
}

// addFlags registers flags for the CLI.
func (o *policyGetOptions) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.myID, "myid", "", "gwctl ID")
}

// run performs the execution of the 'delete policy' subcommand
func (o *policyGetOptions) run() error {
	g, err := admin.GetClientFromID(o.myID)
	if err != nil {
		return err
	}

	rules, err := g.GetACLPolicies()
	if err != nil {
		return err
	}

	fmt.Printf("GW ACL rules\n")
	for r, v := range rules {
		fmt.Printf("[Match]: %v [Action]: %v\n", r, v)
	}

	lPolicies, err := g.GetLBPolicies()
	if err != nil {
		return err
	}

	fmt.Printf("GW Load-balancing policies\n")
	for d, val := range lPolicies {
		for s, p := range val {
			fmt.Printf("ServiceSrc: %v ServiceDst: %v Policy: %v\n", s, d, p)
		}
	}

	return nil
}
