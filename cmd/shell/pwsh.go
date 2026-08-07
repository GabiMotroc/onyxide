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
function o { onyxide open $args }

Set-PSReadLineOption -AddToHistoryHandler {
    param($line)
    $parts = $line -split '\s+'
    onyxide proj add --silent $parts[0] $parts[1]
    $true
}

`
}

func (s *pwsh) Uninit() string {
	return `
Remove-Item function:o -ErrorAction SilentlyContinue
Set-PSReadLineOption -AddToHistoryHandler { param($line) $true }

`
}
