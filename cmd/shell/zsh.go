package shell

type zsh struct{}

func init() {
	Register(&zsh{})
}
func (s *zsh) Name() string {
	return "zsh"
}

func (s *zsh) Init() string {
	return `
alias o='onyxide open'

autoload -Uz add-zsh-hook

_mycli_preexec() {
  typeset -g MYCLI_LAST_CMD="$1"
  typeset -g MYCLI_LAST_PWD="$PWD"
}

_mycli_precmd() {
  if [[ -n "$MYCLI_LAST_CMD" ]]; then
	local parts=("${(@z)MYCLI_LAST_CMD}")
	command onyxide proj add --silent "${parts[1]}" "${parts[2]}" >/dev/null 2>&1 &!
    unset MYCLI_LAST_CMD
    unset MYCLI_LAST_PWD
  fi
}

add-zsh-hook preexec _mycli_preexec
add-zsh-hook precmd _mycli_precmd
`
}

func (s *zsh) Uninit() string {
	return `
add-zsh-hook -d preexec _mycli_preexec
add-zsh-hook -d precmd _mycli_precmd
unfunction _mycli_preexec _mycli_precmd 2>/dev/null
unset MYCLI_LAST_CMD
unalias o 2>/dev/null
`
}
