package command

import (
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
)

func TestGraph(t *testing.T) {
	tmp, cwd := testCwd(t)
	defer testFixCwd(t, tmp, cwd)

	ui := new(cli.MockUi)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
		},
	}

	args := []string{
		testFixturePath("graph"),
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "provider.test") {
		t.Fatalf("doesn't look like digraph: %s", output)
	}
}

func TestGraph_multipleArgs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
		},
	}

	args := []string{
		"bad",
		"bad",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.OutputWriter.String())
	}
}

func TestGraph_noArgs(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := os.Chdir(testFixturePath("graph")); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Chdir(cwd)

	ui := new(cli.MockUi)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "provider.test") {
		t.Fatalf("doesn't look like digraph: %s", output)
	}
}

func TestGraph_noConfig(t *testing.T) {
	td := tempDir(t)
	os.MkdirAll(td, 0755)
	defer os.RemoveAll(td)
	defer testChdir(t, td)()

	ui := new(cli.MockUi)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
		},
	}

	// Running the graph command without a config should not panic,
	// but this may be an error at some point in the future.
	args := []string{"-type", "apply"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
}

func TestGraph_plan(t *testing.T) {
	tmp, cwd := testCwd(t)
	defer testFixCwd(t, tmp, cwd)

	plan := &plans.Plan{
		Changes: plans.NewChanges(),
	}
	plan.Changes.Resources = append(plan.Changes.Resources, &plans.ResourceInstanceChangeSrc{
		Addr: addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "test_instance",
			Name: "bar",
		}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
		ChangeSrc: plans.ChangeSrc{
			Action: plans.Delete,
			Before: plans.DynamicValue(`{}`),
			After:  plans.DynamicValue(`null`),
		},
		ProviderAddr: addrs.ProviderConfig{Type: "test"}.Absolute(addrs.RootModuleInstance),
	})
	emptyConfig, err := plans.NewDynamicValue(cty.EmptyObjectVal, cty.EmptyObject)
	if err != nil {
		t.Fatal(err)
	}
	plan.Backend = plans.Backend{
		// Doesn't actually matter since we aren't going to activate the backend
		// for this command anyway, but we need something here for the plan
		// file writer to succeed.
		Type:   "placeholder",
		Config: emptyConfig,
	}
	_, configSnap := testModuleWithSnapshot(t, "graph")

	planPath := testPlanFile(t, configSnap, states.NewState(), plan)

	ui := new(cli.MockUi)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
		},
	}

	args := []string{
		planPath,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "provider.test") {
		t.Fatalf("doesn't look like digraph: %s", output)
	}
}
