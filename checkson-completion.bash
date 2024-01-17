# bash completion for checkson                             -*- shell-script -*-

__checkson_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__checkson_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__checkson_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__checkson_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__checkson_handle_go_custom_completion()
{
    __checkson_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=1
    local shellCompDirectiveNoSpace=2
    local shellCompDirectiveNoFileComp=4
    local shellCompDirectiveFilterFileExt=8
    local shellCompDirectiveFilterDirs=16

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly checkson allows to handle aliases
    args=("${words[@]:1}")
    requestComp="${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __checkson_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __checkson_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __checkson_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __checkson_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __checkson_debug "${FUNCNAME[0]}: the completions are: ${out[*]}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __checkson_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __checkson_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __checkson_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out[*]}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __checkson_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subDir
        # Use printf to strip any trailing newline
        subdir=$(printf "%s" "${out[0]}")
        if [ -n "$subdir" ]; then
            __checkson_debug "Listing directories in $subdir"
            __checkson_handle_subdirs_in_dir_flag "$subdir"
        else
            __checkson_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out[*]}" -- "$cur")
    fi
}

__checkson_handle_reply()
{
    __checkson_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __checkson_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __checkson_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __checkson_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		if declare -F __checkson_custom_func >/dev/null; then
			# try command name qualified custom func
			__checkson_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__checkson_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__checkson_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__checkson_handle_flag()
{
    __checkson_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __checkson_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __checkson_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __checkson_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __checkson_contains_word "${words[c]}" "${two_word_flags[@]}"; then
			  __checkson_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__checkson_handle_noun()
{
    __checkson_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __checkson_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __checkson_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__checkson_handle_command()
{
    __checkson_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_checkson_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __checkson_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__checkson_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __checkson_handle_reply
        return
    fi
    __checkson_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __checkson_handle_flag
    elif __checkson_contains_word "${words[c]}" "${commands[@]}"; then
        __checkson_handle_command
    elif [[ $c -eq 0 ]]; then
        __checkson_handle_command
    elif __checkson_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __checkson_handle_command
        else
            __checkson_handle_noun
        fi
    else
        __checkson_handle_noun
    fi
    __checkson_handle_word
}

_checkson_channels_create()
{
    last_command="checkson_channels_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--email=")
    two_word_flags+=("--email")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--email=")
    flags+=("--pager-duty-service-key=")
    two_word_flags+=("--pager-duty-service-key")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--pager-duty-service-key=")
    flags+=("--slack-incoming-webhook-url=")
    two_word_flags+=("--slack-incoming-webhook-url")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--slack-incoming-webhook-url=")
    flags+=("--type=")
    two_word_flags+=("--type")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--type=")
    flags+=("--webhook-url=")
    two_word_flags+=("--webhook-url")
    two_word_flags+=("-w")
    local_nonpersistent_flags+=("--webhook-url=")
    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_channels_delete()
{
    last_command="checkson_channels_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_channels_list()
{
    last_command="checkson_channels_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_channels_show()
{
    last_command="checkson_channels_show"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_channels()
{
    last_command="checkson_channels"

    command_aliases=()

    commands=()
    commands+=("create")
    commands+=("delete")
    commands+=("list")
    commands+=("show")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_completion()
{
    last_command="checkson_completion"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("bash")
    must_have_one_noun+=("fish")
    must_have_one_noun+=("powershell")
    must_have_one_noun+=("zsh")
    noun_aliases=()
}

_checkson_create()
{
    last_command="checkson_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--channel=")
    two_word_flags+=("--channel")
    two_word_flags+=("-c")
    local_nonpersistent_flags+=("--channel=")
    flags+=("--check-interval=")
    two_word_flags+=("--check-interval")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--check-interval=")
    flags+=("--docker-image=")
    two_word_flags+=("--docker-image")
    two_word_flags+=("-i")
    local_nonpersistent_flags+=("--docker-image=")
    flags+=("--email=")
    two_word_flags+=("--email")
    two_word_flags+=("-m")
    local_nonpersistent_flags+=("--email=")
    flags+=("--env=")
    two_word_flags+=("--env")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--env=")
    flags+=("--failure-threshold=")
    two_word_flags+=("--failure-threshold")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--failure-threshold=")
    flags+=("--webhook-url=")
    two_word_flags+=("--webhook-url")
    two_word_flags+=("-w")
    local_nonpersistent_flags+=("--webhook-url=")
    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_delete()
{
    last_command="checkson_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_help()
{
    last_command="checkson_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_checkson_list()
{
    last_command="checkson_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_login()
{
    last_command="checkson_login"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--personal-access-token=")
    two_word_flags+=("--personal-access-token")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--personal-access-token=")
    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_logout()
{
    last_command="checkson_logout"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_logs()
{
    last_command="checkson_logs"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_runs()
{
    last_command="checkson_runs"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_show()
{
    last_command="checkson_show"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_status()
{
    last_command="checkson_status"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_version()
{
    last_command="checkson_version"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_checkson_root_command()
{
    last_command="checkson"

    command_aliases=()

    commands=()
    commands+=("channels")
    commands+=("completion")
    commands+=("create")
    commands+=("delete")
    commands+=("help")
    commands+=("list")
    commands+=("login")
    commands+=("logout")
    commands+=("logs")
    commands+=("runs")
    commands+=("show")
    commands+=("status")
    commands+=("version")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config-file=")
    two_word_flags+=("--config-file")
    two_word_flags+=("-C")
    flags+=("--dev-mode")
    flags+=("--verbose")
    flags+=("-V")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_checkson()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __checkson_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("checkson")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function
    local last_command
    local nouns=()

    __checkson_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_checkson checkson
else
    complete -o default -o nospace -F __start_checkson checkson
fi

# ex: ts=4 sw=4 et filetype=sh
