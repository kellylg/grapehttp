package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"grapehttp/client/cmd/templates"
	cmdutil "grapehttp/client/cmd/util"
)

const (
	bashCompletionFunc = `# call fctl get $1,
__fctl_override_flag_list=(gconfig cluster user context namespace server)
__fctl_override_flags()
{
    local ${__fctl_override_flag_list[*]} two_word_of of
    for w in "${words[@]}"; do
        if [ -n "${two_word_of}" ]; then
            eval "${two_word_of}=\"--${two_word_of}=\${w}\""
            two_word_of=
            continue
        fi
        for of in "${__fctl_override_flag_list[@]}"; do
            case "${w}" in
                --${of}=*)
                    eval "${of}=\"${w}\""
                    ;;
                --${of})
                    two_word_of="${of}"
                    ;;
            esac
        done
        if [ "${w}" == "--all-namespaces" ]; then
            namespace="--all-namespaces"
        fi
    done
    for of in "${__fctl_override_flag_list[@]}"; do
        if eval "test -n \"\$${of}\""; then
            eval "echo \${${of}}"
        fi
    done
}

__fctl_get_namespaces()
{
    local template fctl_out
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    if fctl_out=$(fctl get -o template --template="${template}" namespace 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${fctl_out[*]}" -- "$cur" ) )
    fi
}

__fctl_config_get_contexts()
{
    __fctl_parse_config "contexts"
}

__fctl_config_get_clusters()
{
    __fctl_parse_config "clusters"
}

__fctl_config_get_users()
{
    __fctl_parse_config "users"
}

# $1 has to be "contexts", "clusters" or "users"
__fctl_config_get()
{
    local template fctl_out
    template="{{ range .$1  }}{{ .name }} {{ end }}"
    if fctl_out=$(fctl config $(__fctl_override_flags) -o template --template="${template}" view 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${fctl_out[*]}" -- "$cur" ) )
    fi
}

__fctl_parse_get()
{
    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local fctl_out
    if fctl_out=$(fctl get $(__fctl_override_flags) -o template --template="${template}" "$1" 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${fctl_out[*]}" -- "$cur" ) )
    fi
}

__fctl_get_resource()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        return 1
    fi
    __fctl_parse_get "${nouns[${#nouns[@]} -1]}"
}

__fctl_get_resource_pod()
{
    __fctl_parse_get "pod"
}

__fctl_get_resource_rc()
{
    __fctl_parse_get "rc"
}

__fctl_get_resource_node()
{
    __fctl_parse_get "node"
}

# $1 is the name of the pod we want to get the list of containers inside
__fctl_get_containers()
{
    local template
    template="{{ range .spec.containers  }}{{ .name }} {{ end }}"
    __debug "${FUNCNAME} nouns are ${nouns[*]}"

    local len="${#nouns[@]}"
    if [[ ${len} -ne 1 ]]; then
        return
    fi
    local last=${nouns[${len} -1]}
    local fctl_out
    if fctl_out=$(fctl get $(__fctl_override_flags) -o template --template="${template}" pods "${last}" 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${fctl_out[*]}" -- "$cur" ) )
    fi
}

# Require both a pod and a container to be specified
__fctl_require_pod_and_container()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        __fctl_parse_get pods
        return 0
    fi;
    __fctl_get_containers
    return 0
}

__custom_func() {
    case ${last_command} in
        fctl_get | fctl_describe | fctl_delete | fctl_label | fctl_stop | fctl_edit | fctl_patch |\
        fctl_annotate | fctl_expose | fctl_scale | fctl_autoscale | fctl_taint | fctl_rollout_*)
            __fctl_get_resource
            return
            ;;
        fctl_logs | fctl_attach)
            __fctl_require_pod_and_container
            return
            ;;
        fctl_exec | fctl_port-forward | fctl_top_pod)
            __fctl_get_resource_pod
            return
            ;;
        fctl_rolling-update)
            __fctl_get_resource_rc
            return
            ;;
        fctl_cordon | fctl_uncordon | fctl_drain | fctl_top_node)
            __fctl_get_resource_node
            return
            ;;
        fctl_config_use-context)
            __fctl_config_get_contexts
            return
            ;;
        *)
            ;;
    esac
}
`
)

var (
	bash_completion_flags = map[string]string{
		"namespace": "__fctl_get_namespaces",
		"context":   "__fctl_config_get_contexts",
		"cluster":   "__fctl_config_get_clusters",
		"user":      "__fctl_config_get_users",
	}
)

func NewFctlCommand(f cmdutil.Factory, in io.Reader, out, err io.Writer) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "fctl",
		Short: "grapehttp server cli tool",
		Long: templates.LongDesc(`
      fctl is used to mkdir, ls, cp, mv, rm, download, upload files in the http file server.`),
		Run: runHelp,
		BashCompletionFunction: bashCompletionFunc,
	}

	groups := templates.CommandGroups{
		{
			Message: "File Server Commands:",
			Commands: []*cobra.Command{
				NewCmdLs(f, out, err),
				NewCmdMkdir(f, out, err),
				NewCmdMv(f, out, err),
				NewCmdCp(f, out, err),
				NewCmdRm(f, out, err),
				NewCmdUpload(f, out, err),
				NewCmdDownload(f, out, err),
			},
		},
		{
			Message: "User Control Commands:",
			Commands: []*cobra.Command{
				NewCmdUserAdd(f, out, err),
				NewCmdUserDel(f, out, err),
				NewCmdUserModify(f, out, err),
				NewCmdUserSearch(f, out, err),
				NewCmdUserGet(f, out, err),
				NewCmdUserList(f, out, err),
				NewCmdUserEnable(f, out, err),
				NewCmdUserDisable(f, out, err),
			},
		},
	}
	groups.Add(cmds)
	templates.ActsAsRootCommand(cmds, []string{}, groups...)

	for name, completion := range bash_completion_flags {
		if cmds.Flag(name) != nil {
			if cmds.Flag(name).Annotations == nil {
				cmds.Flag(name).Annotations = map[string][]string{}
			}
			cmds.Flag(name).Annotations[cobra.BashCompCustom] = append(
				cmds.Flag(name).Annotations[cobra.BashCompCustom],
				completion,
			)
		}
	}

	cmds.AddCommand(NewCmdVersion(f, out))
	cmds.AddCommand(NewCmdCompletion(out, ""))
	//cmds.AddCommand(NewCmdOptions(out))
	cmds.AddCommand(NewCmdFinfo(f, out, err))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
