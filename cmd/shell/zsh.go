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
alias o='onyxide'

autoload -Uz add-zsh-hook

_mycli_preexec() {
  typeset -g MYCLI_LAST_CMD="$1"
  typeset -g MYCLI_LAST_PWD="$PWD"
}

_mycli_precmd() {
  if [[ -n "$MYCLI_LAST_CMD" ]]; then
    command onyxide hook --pwd "$MYCLI_LAST_PWD" --raw "$MYCLI_LAST_CMD"
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
`
}
