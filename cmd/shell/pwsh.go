package shell

type pwsh struct{}

func init() {
	Register(&pwsh{})
}
func (s *pwsh) Name() string {
	return "pwsh"
}

func (s *pwsh) Init() string {
	return `
Set-Alias -Name o -Value onyxide

Set-PSReadLineOption -AddToHistoryHandler {
    param($line)
    onyxide hook --pwd $PWD.Path --raw $line
    $true
}
`
}

func (s *pwsh) Uninit() string {
	return `
Remove-Alias -Name o -ErrorAction SilentlyContinue
Set-PSReadLineOption -AddToHistoryHandler { param($line) $true }

`
}
