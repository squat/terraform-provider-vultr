package configupgrade

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/zclconf/go-cty/cty"

	backendinit "github.com/hashicorp/terraform/backend/init"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/provisioners"
	"github.com/hashicorp/terraform/terraform"
)

func TestUpgradeValid(t *testing.T) {
	// This test uses the contents of the test-fixtures/valid directory as
	// a table of tests. Every directory there must have both "input" and
	// "want" subdirectories, where "input" is the configuration to be
	// upgraded and "want" is the expected result.
	fixtureDir := "test-fixtures/valid"
	testDirs, err := ioutil.ReadDir(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range testDirs {
		if !entry.IsDir() {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			inputDir := filepath.Join(fixtureDir, entry.Name(), "input")
			wantDir := filepath.Join(fixtureDir, entry.Name(), "want")
			u := &Upgrader{
				Providers:    providers.ResolverFixed(testProviders),
				Provisioners: testProvisioners,
			}

			inputSrc, err := LoadModule(inputDir)
			if err != nil {
				t.Fatal(err)
			}
			wantSrc, err := LoadModule(wantDir)
			if err != nil {
				t.Fatal(err)
			}

			gotSrc, diags := u.Upgrade(inputSrc, inputDir)
			if diags.HasErrors() {
				t.Error(diags.Err())
			}

			// Upgrade uses a nil entry as a signal to delete a file, which
			// we can't test here because we aren't modifying an existing
			// dir in place, so we'll just ignore those and leave that mechanism
			// to be tested elsewhere.

			for name, got := range gotSrc {
				if gotSrc[name] == nil {
					delete(gotSrc, name)
					continue
				}
				want, wanted := wantSrc[name]
				if !wanted {
					t.Errorf("unexpected extra output file %q\n=== GOT ===\n%s", name, got)
					continue
				}

				got = bytes.TrimSpace(got)
				want = bytes.TrimSpace(want)
				if !bytes.Equal(got, want) {
					diff := diffSourceFiles(got, want)
					t.Errorf("wrong content in %q\n%s", name, diff)
				}
			}

			for name, want := range wantSrc {
				if _, present := gotSrc[name]; !present {
					t.Errorf("missing output file %q\n=== WANT ===\n%s", name, want)
				}
			}
		})
	}
}

func TestUpgradeRenameJSON(t *testing.T) {
	inputDir := filepath.Join("test-fixtures/valid/rename-json/input")
	inputSrc, err := LoadModule(inputDir)
	if err != nil {
		t.Fatal(err)
	}

	u := &Upgrader{
		Providers: providers.ResolverFixed(testProviders),
	}
	gotSrc, diags := u.Upgrade(inputSrc, inputDir)
	if diags.HasErrors() {
		t.Error(diags.Err())
	}

	// This test fixture is also fully covered by TestUpgradeValid, so
	// we're just testing that the file was renamed here.
	src, exists := gotSrc["misnamed-json.tf"]
	if src != nil {
		t.Errorf("misnamed-json.tf still has content")
	} else if !exists {
		t.Errorf("misnamed-json.tf not marked for deletion")
	}

	src, exists = gotSrc["misnamed-json.tf.json"]
	if src == nil || !exists {
		t.Errorf("misnamed-json.tf.json was not created")
	}
}

func diffSourceFiles(got, want []byte) []byte {
	// We'll try to run "diff -u" here to get nice output, but if that fails
	// (e.g. because we're running on a machine without diff installed) then
	// we'll fall back on just printing out the before and after in full.
	gotR, gotW, err := os.Pipe()
	if err != nil {
		return diffSourceFilesFallback(got, want)
	}
	defer gotR.Close()
	defer gotW.Close()
	wantR, wantW, err := os.Pipe()
	if err != nil {
		return diffSourceFilesFallback(got, want)
	}
	defer wantR.Close()
	defer wantW.Close()

	cmd := exec.Command("diff", "-u", "--label=GOT", "--label=WANT", "/dev/fd/3", "/dev/fd/4")
	cmd.ExtraFiles = []*os.File{gotR, wantR}
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return diffSourceFilesFallback(got, want)
	}

	go func() {
		wantW.Write(want)
		wantW.Close()
	}()
	go func() {
		gotW.Write(got)
		gotW.Close()
	}()

	err = cmd.Start()
	if err != nil {
		return diffSourceFilesFallback(got, want)
	}

	outR := io.MultiReader(stdout, stderr)
	out, err := ioutil.ReadAll(outR)
	if err != nil {
		return diffSourceFilesFallback(got, want)
	}

	cmd.Wait() // not checking errors here because on failure we'll have stderr captured to return

	const noNewline = "\\ No newline at end of file\n"
	if bytes.HasSuffix(out, []byte(noNewline)) {
		out = out[:len(out)-len(noNewline)]
	}
	return out
}

func diffSourceFilesFallback(got, want []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString("=== GOT ===\n")
	buf.Write(got)
	buf.WriteString("\n=== WANT ===\n")
	buf.Write(want)
	buf.WriteString("\n")
	return buf.Bytes()
}

var testProviders = map[string]providers.Factory{
	"test": providers.Factory(func() (providers.Interface, error) {
		p := &terraform.MockProvider{}
		p.GetSchemaReturn = &terraform.ProviderSchema{
			ResourceTypes: map[string]*configschema.Block{
				"test_instance": {
					Attributes: map[string]*configschema.Attribute{
						"id":              {Type: cty.String, Computed: true},
						"type":            {Type: cty.String, Optional: true},
						"image":           {Type: cty.String, Optional: true},
						"tags":            {Type: cty.Map(cty.String), Optional: true},
						"security_groups": {Type: cty.List(cty.String), Optional: true},
						"subnet_ids":      {Type: cty.Set(cty.String), Optional: true},
						"list_of_obj":     {Type: cty.List(cty.EmptyObject), Optional: true},
					},
					BlockTypes: map[string]*configschema.NestedBlock{
						"network": {
							Nesting: configschema.NestingSet,
							Block: configschema.Block{
								Attributes: map[string]*configschema.Attribute{
									"cidr_block":   {Type: cty.String, Optional: true},
									"subnet_cidrs": {Type: cty.Map(cty.String), Computed: true},
								},
								BlockTypes: map[string]*configschema.NestedBlock{
									"subnet": {
										Nesting: configschema.NestingSet,
										Block: configschema.Block{
											Attributes: map[string]*configschema.Attribute{
												"number": {Type: cty.Number, Required: true},
											},
										},
									},
								},
							},
						},
						"addresses": {
							Nesting: configschema.NestingSingle,
							Block: configschema.Block{
								Attributes: map[string]*configschema.Attribute{
									"ipv4": {Type: cty.String, Computed: true},
									"ipv6": {Type: cty.String, Computed: true},
								},
							},
						},
					},
				},
			},
		}
		return p, nil
	}),
	"terraform": providers.Factory(func() (providers.Interface, error) {
		p := &terraform.MockProvider{}
		p.GetSchemaReturn = &terraform.ProviderSchema{
			DataSources: map[string]*configschema.Block{
				"terraform_remote_state": {
					// This is just enough an approximation of the remote state
					// schema to check out reference upgrade logic. It is
					// intentionally not fully-comprehensive.
					Attributes: map[string]*configschema.Attribute{
						"backend": {Type: cty.String, Optional: true},
					},
				},
			},
		}
		return p, nil
	}),
	"aws": providers.Factory(func() (providers.Interface, error) {
		// This is here only so we can test the provisioner connection info
		// migration behavior, which is resource-type specific. Do not use
		// it in any other tests.
		p := &terraform.MockProvider{}
		p.GetSchemaReturn = &terraform.ProviderSchema{
			ResourceTypes: map[string]*configschema.Block{
				"aws_instance": {},
			},
		}
		return p, nil
	}),
}

var testProvisioners = map[string]provisioners.Factory{
	"test": provisioners.Factory(func() (provisioners.Interface, error) {
		p := &terraform.MockProvisioner{}
		p.GetSchemaResponse = provisioners.GetSchemaResponse{
			Provisioner: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"commands":    {Type: cty.List(cty.String), Optional: true},
					"interpreter": {Type: cty.String, Optional: true},
				},
			},
		}
		return p, nil
	}),
}

func init() {
	// Initialize the backends
	backendinit.Init(nil)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		// if we're verbose, use the logging requested by TF_LOG
		logging.SetOutput()
	} else {
		// otherwise silence all logs
		log.SetOutput(ioutil.Discard)
	}

	// We have fmt.Stringer implementations on lots of objects that hide
	// details that we very often want to see in tests, so we just disable
	// spew's use of String methods globally on the assumption that spew
	// usage implies an intent to see the raw values and ignore any
	// abstractions.
	spew.Config.DisableMethods = true

	os.Exit(m.Run())
}
