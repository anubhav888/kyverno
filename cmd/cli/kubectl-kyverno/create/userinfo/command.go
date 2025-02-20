package userinfo

import (
	"os"
	"text/template"

	"github.com/kyverno/kyverno/api/kyverno/v1beta1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/create/templates"
	"github.com/spf13/cobra"
	authenticationv1 "k8s.io/api/authentication/v1"
)

func Command() *cobra.Command {
	var path string
	var username string
	var roles, clusterRoles, groups []string
	cmd := &cobra.Command{
		Use:     "user-info",
		Short:   "Create a Kyverno user-info file.",
		Example: "kyverno create user-info -u molybdenum@somecorp.com -g basic-user -c admin",
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpl, err := template.New("userinfo").Parse(templates.UserInfoTemplate)
			if err != nil {
				return err
			}
			output := os.Stdout
			if path != "" {
				file, err := os.Create(path)
				if err != nil {
					return err
				}
				defer file.Close()
				output = file
			}
			values := v1beta1.RequestInfo{
				Roles:        roles,
				ClusterRoles: clusterRoles,
				AdmissionUserInfo: authenticationv1.UserInfo{
					Username: username,
					Groups:   groups,
				},
			}
			return tmpl.Execute(output, values)
		},
	}
	cmd.Flags().StringVarP(&path, "output", "o", "", "Output path (uses standard console output if not set)")
	cmd.Flags().StringVarP(&username, "username", "u", "", "User name")
	cmd.Flags().StringSliceVarP(&roles, "role", "r", nil, "Role")
	cmd.Flags().StringSliceVarP(&clusterRoles, "cluster-role", "c", nil, "Cluster role")
	cmd.Flags().StringSliceVarP(&groups, "group", "g", nil, "Group")
	return cmd
}
