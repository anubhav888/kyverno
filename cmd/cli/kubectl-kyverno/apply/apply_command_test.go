package apply

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	policyreportv1alpha2 "github.com/kyverno/kyverno/api/policyreport/v1alpha2"
	reportutils "github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/report"
	"gotest.tools/assert"
)

func Test_Apply(t *testing.T) {
	type TestCase struct {
		gitBranch             string
		expectedPolicyReports []policyreportv1alpha2.PolicyReport
		config                ApplyCommandConfig
		stdinFile             string
	}
	// copy disallow_latest_tag.yaml to local path
	localFileName, err := copyFileToThisDir("../../../../test/best_practices/disallow_latest_tag.yaml")
	assert.NilError(t, err)
	defer func() { _ = os.Remove(localFileName) }()

	testcases := []*TestCase{{
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/best_practices/disallow_latest_tag.yaml"},
			ResourcePaths: []string{"../../../../test/resources/pod_with_version_tag.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  2,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{localFileName},
			ResourcePaths: []string{"../../../../test/resources/pod_with_version_tag.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  2,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/best_practices/disallow_latest_tag.yaml"},
			ResourcePaths: []string{"../../../../test/resources/pod_with_latest_tag.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  1,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/apply/policies"},
			ResourcePaths: []string{"../../../../test/cli/apply/resource"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  1,
				Skip:  8,
				Error: 0,
				Warn:  2,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/best_practices/disallow_latest_tag.yaml"},
			ResourcePaths: []string{"../../../../test/resources/pod_with_latest_tag.yaml"},
			PolicyReport:  true,
			AuditWarn:     true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  1,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"-"},
			ResourcePaths: []string{"../../../../test/resources/pod_with_latest_tag.yaml"},
			PolicyReport:  true,
			AuditWarn:     true,
		},
		stdinFile: "../../../../test/best_practices/disallow_latest_tag.yaml",
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  1,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/best_practices/disallow_latest_tag.yaml"},
			ResourcePaths: []string{"-"},
			PolicyReport:  true,
			AuditWarn:     true,
		},
		stdinFile: "../../../../test/resources/pod_with_latest_tag.yaml",
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  1,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"https://github.com/kyverno/policies/openshift/team-validate-ns-name/"},
			ResourcePaths: []string{"../../../../test/openshift/team-validate-ns-name.yaml"},
			GitBranch:     "main",
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  2,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/apply/policies-set"},
			ResourcePaths: []string{"../../../../test/cli/apply/resources-set"},
			Variables:     []string{"request.operation=UPDATE"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  2,
				Fail:  0,
				Skip:  4,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/test-validating-admission-policy/check-deployments-replica/policy.yaml"},
			ResourcePaths: []string{"../../../../test/cli/test-validating-admission-policy/check-deployments-replica/deployment1.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/test-validating-admission-policy/check-deployments-replica/policy.yaml"},
			ResourcePaths: []string{"../../../../test/cli/test-validating-admission-policy/check-deployments-replica/deployment2.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  0,
				Fail:  1,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/test-validating-admission-policy/disallow-host-path/policy.yaml"},
			ResourcePaths: []string{"../../../../test/cli/test-validating-admission-policy/disallow-host-path/pod1.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/test-validating-admission-policy/disallow-host-path/policy.yaml"},
			ResourcePaths: []string{"../../../../test/cli/test-validating-admission-policy/disallow-host-path/pod2.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  0,
				Fail:  1,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/test-validating-admission-policy/check-deployment-labels/policy.yaml"},
			ResourcePaths: []string{"../../../../test/cli/test-validating-admission-policy/check-deployment-labels/deployment1.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  1,
				Fail:  0,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/cli/test-validating-admission-policy/check-deployment-labels/policy.yaml"},
			ResourcePaths: []string{"../../../../test/cli/test-validating-admission-policy/check-deployment-labels/deployment2.yaml"},
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  0,
				Fail:  1,
				Skip:  0,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"https://github.com/kyverno/policies/best-practices/require-labels/", "../../../../test/best_practices/disallow_latest_tag.yaml"},
			ResourcePaths: []string{"../../../../test/resources/pod_with_version_tag.yaml"},
			GitBranch:     "main",
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  2,
				Fail:  1,
				Skip:  2,
				Error: 0,
				Warn:  0,
			},
		}},
	}, {
		// Same as the above test case but the policy paths are reordered
		config: ApplyCommandConfig{
			PolicyPaths:   []string{"../../../../test/best_practices/disallow_latest_tag.yaml", "https://github.com/kyverno/policies/best-practices/require-labels/"},
			ResourcePaths: []string{"../../../../test/resources/pod_with_version_tag.yaml"},
			GitBranch:     "main",
			PolicyReport:  true,
		},
		expectedPolicyReports: []policyreportv1alpha2.PolicyReport{{
			Summary: policyreportv1alpha2.PolicyReportSummary{
				Pass:  2,
				Fail:  1,
				Skip:  2,
				Error: 0,
				Warn:  0,
			},
		}},
	}}

	compareSummary := func(expected policyreportv1alpha2.PolicyReportSummary, actual policyreportv1alpha2.PolicyReportSummary, desc string) {
		assert.Equal(t, actual.Pass, expected.Pass, desc)
		assert.Equal(t, actual.Fail, expected.Fail, desc)
		assert.Equal(t, actual.Skip, expected.Skip, desc)
		assert.Equal(t, actual.Warn, expected.Warn, desc)
		assert.Equal(t, actual.Error, expected.Error, desc)
	}

	verifyTestcase := func(t *testing.T, tc *TestCase, compareSummary func(policyreportv1alpha2.PolicyReportSummary, policyreportv1alpha2.PolicyReportSummary, string)) {
		if tc.stdinFile != "" {
			oldStdin := os.Stdin
			input, err := os.OpenFile(tc.stdinFile, os.O_RDONLY, 0)
			assert.NilError(t, err)
			os.Stdin = input
			defer func() {
				// Restore original Stdin
				os.Stdin = oldStdin
				_ = input.Close()
			}()
		}
		desc := fmt.Sprintf("Policies: [%s], / Resources: [%s]", strings.Join(tc.config.PolicyPaths, ","), strings.Join(tc.config.ResourcePaths, ","))
		// prevent os.Exit from being called
		osExit = func(code int) {
			assert.Check(t, false, "os.Exit(%d) should not be called: %s", code, desc)
		}

		defer func() { osExit = os.Exit }()
		_, _, _, responses, err := tc.config.applyCommandHelper()
		assert.NilError(t, err, desc)

		clustered, _ := reportutils.ComputePolicyReports(tc.config.AuditWarn, responses...)
		assert.Assert(t, len(clustered) > 0, "policy reports should not be empty: %s", desc)
		combined := []policyreportv1alpha2.ClusterPolicyReport{
			reportutils.MergeClusterReports(clustered, nil),
		}
		assert.Equal(t, len(combined), len(tc.expectedPolicyReports))
		for i, resp := range combined {
			compareSummary(tc.expectedPolicyReports[i].Summary, resp.Summary, desc)
		}
	}

	for _, tc := range testcases {
		verifyTestcase(t, tc, compareSummary)
	}
}

func copyFileToThisDir(sourceFile string) (string, error) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", err
	}

	return filepath.Base(sourceFile), os.WriteFile(filepath.Base(sourceFile), input, 0o644)
}
