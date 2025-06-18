package browser

var (
	Reset         string
	Bold          string
	Italic        string
	FgWhite       string
	FgCyan        string
	FgPink        string
	FgRed         string
	FgGray        string
	FgYellow      string
	FgGreen       string
	BgGreen       string
	BgPurple      string
	VimSelectBg   string
	VimSelectFg   string
	FgWhiteBright string
)

func init() {
	Reset = "\033[0m"
	Bold = "\033[1m"
	Italic = "\033[3m"
	FgWhite = "\033[97m"
	FgCyan = "\033[96m"
	FgPink = "\033[95m"
	FgRed = "\033[91m"
	FgGray = "\033[100m"
	FgYellow = "\033[93m"
	FgGreen = "\033[42m"
	BgPurple = "\033[45m"
	VimSelectBg = "\033[43m"
	VimSelectFg = "\033[93m"
	BgGreen = "\033[42m"
	FgWhiteBright = "\033[97m"
}
