package format

import (
	"testing"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/colorstring"
	"github.com/zclconf/go-cty/cty"
)

var disabledColorize = &colorstring.Colorize{
	Colors:  colorstring.DefaultColors,
	Disable: true,
}

func TestState(t *testing.T) {
	state := states.NewState()

	rootModule := state.RootModule()
	if rootModule == nil {
		t.Errorf("root module is nil; want valid object")
	}

	rootModule.SetLocalValue("foo", cty.StringVal("foo value"))
	rootModule.SetOutputValue("bar", cty.StringVal("bar value"), false)
	rootModule.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "test_resource",
			Name: "baz",
		}.Instance(addrs.IntKey(0)),
		&states.ResourceInstanceObjectSrc{
			Status:        states.ObjectReady,
			SchemaVersion: 1,
			AttrsJSON:     []byte(`{"woozles":"confuzles"}`),
		},
		addrs.ProviderConfig{
			Type: "test",
		}.Absolute(addrs.RootModuleInstance),
	)
	rootModule.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.DataResourceMode,
			Type: "test_data_source",
			Name: "data",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:        states.ObjectReady,
			SchemaVersion: 1,
			AttrsJSON:     []byte(`{"compute":"sure"}`),
		},
		addrs.ProviderConfig{
			Type: "test",
		}.Absolute(addrs.RootModuleInstance),
	)

	tests := []struct {
		State *StateOpts
		Want  string
	}{
		{
			&StateOpts{
				State:   &states.State{},
				Color:   disabledColorize,
				Schemas: &terraform.Schemas{},
			},
			"The state file is empty. No resources are represented.",
		},
		{
			&StateOpts{
				State:   state,
				Color:   disabledColorize,
				Schemas: testSchemas(),
			},
			TestOutput,
		},
		{
			&StateOpts{
				State:   nestedState(t),
				Color:   disabledColorize,
				Schemas: testSchemas(),
			},
			nestedTestOutput,
		},
	}

	for _, tt := range tests {
		got := State(tt.State)
		if got != tt.Want {
			t.Errorf(
				"wrong result\ninput: %v\ngot: \n%s\nwant: \n%s",
				tt.State.State, got, tt.Want,
			)
		}
	}
}

func testProvider() *terraform.MockProvider {
	p := new(terraform.MockProvider)
	p.ReadResourceFn = func(req providers.ReadResourceRequest) providers.ReadResourceResponse {
		return providers.ReadResourceResponse{NewState: req.PriorState}
	}

	p.GetSchemaReturn = testProviderSchema()

	return p
}

func testProviderSchema() *terraform.ProviderSchema {
	return &terraform.ProviderSchema{
		Provider: &configschema.Block{
			Attributes: map[string]*configschema.Attribute{
				"region": {Type: cty.String, Optional: true},
			},
		},
		ResourceTypes: map[string]*configschema.Block{
			"test_resource": {
				Attributes: map[string]*configschema.Attribute{
					"id":      {Type: cty.String, Computed: true},
					"foo":     {Type: cty.String, Optional: true},
					"woozles": {Type: cty.String, Optional: true},
				},
				BlockTypes: map[string]*configschema.NestedBlock{
					"nested": {
						Nesting: configschema.NestingList,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"compute": {Type: cty.String, Optional: true},
								"value":   {Type: cty.String, Optional: true},
							},
						},
					},
				},
			},
		},
		DataSources: map[string]*configschema.Block{
			"test_data_source": {
				Attributes: map[string]*configschema.Attribute{
					"compute": {Type: cty.String, Optional: true},
					"value":   {Type: cty.String, Computed: true},
				},
			},
		},
	}
}

func testSchemas() *terraform.Schemas {
	provider := testProvider()
	return &terraform.Schemas{
		Providers: map[string]*terraform.ProviderSchema{
			"test": provider.GetSchemaReturn,
		},
	}
}

const TestOutput = `# data.test_data_source.data: 
data "test_data_source" "data" {
    compute = "sure"
}

# test_resource.baz[0]: 
resource "test_resource" "baz" {
    woozles = "confuzles"
}


Outputs:

bar = "bar value"`

const nestedTestOutput = `# test_resource.baz[0]: 
resource "test_resource" "baz" {
    woozles = "confuzles"

    nested {
        value = "42"
    }
}

`

func nestedState(t *testing.T) *states.State {
	state := states.NewState()

	rootModule := state.RootModule()
	if rootModule == nil {
		t.Errorf("root module is nil; want valid object")
	}

	rootModule.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "test_resource",
			Name: "baz",
		}.Instance(addrs.IntKey(0)),
		&states.ResourceInstanceObjectSrc{
			Status:        states.ObjectReady,
			SchemaVersion: 1,
			AttrsJSON:     []byte(`{"woozles":"confuzles","nested": [{"value": "42"}]}`),
		},
		addrs.ProviderConfig{
			Type: "test",
		}.Absolute(addrs.RootModuleInstance),
	)
	return state
}
