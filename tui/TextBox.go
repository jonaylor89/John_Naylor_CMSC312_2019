package tui

import (
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
)

const (
	tokenFg       = "fg"
	tokenBg       = "bg"
	tokenModifier = "mod"

	tokenItemSeparator  = ","
	tokenValueSeparator = ":"

	tokenBeginStyledText = '['
	tokenEndStyledText   = ']'

	tokenBeginStyle = '('
	tokenEndStyle   = ')'
)

var (
	textColorMap = map[Color]string{
		ColorRed:     "red",
		ColorBlue:    "blue",
		ColorBlack:   "black",
		ColorCyan:    "cyan",
		ColorYellow:  "yellow",
		ColorWhite:   "white",
		ColorClear:   "clear",
		ColorGreen:   "green",
		ColorMagenta: "magenta",
	}

	textModifierMap = map[Modifier]string{
		ModifierBold:      "bold",
		ModifierUnderline: "underline",
		ModifierReverse:   "reverse",
	}
)

type TextBox struct {
	Block
	WrapText    bool
	TextStyle   Style
	CursorStyle Style
	ShowCursor  bool

	text        [][]Cell
	cursorPoint image.Point
}

func NewTextBox() *TextBox {
	return &TextBox{
		Block:       *NewBlock(),
		WrapText:    false,
		TextStyle:   NewStyle(ColorWhite),
		CursorStyle: NewStyle(ColorWhite, ColorClear, ModifierReverse),

		text:        [][]Cell{[]Cell{}},
		cursorPoint: image.Pt(1, 1),
	}
}

func (self *TextBox) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	yCoordinate := 0
	for _, line := range self.text {
		if self.WrapText {
			line = WrapCells(line, uint(self.Inner.Dx()))
		}
		lines := SplitCells(line, '\n')
		for _, line := range lines {
			for _, cx := range BuildCellWithXArray(line) {
				x, cell := cx.X, cx.Cell
				buf.SetCell(cell, image.Pt(x, yCoordinate).Add(self.Inner.Min))
			}
			yCoordinate++
		}
		if yCoordinate > self.Inner.Max.Y {
			break
		}
	}

	if self.ShowCursor {
		point := self.cursorPoint.Add(self.Inner.Min).Sub(image.Pt(1, 1))
		cell := buf.GetCell(point)
		cell.Style = self.CursorStyle
		buf.SetCell(cell, point)
	}
}

func (self *TextBox) Backspace() {
	if self.cursorPoint == image.Pt(17, 1) {
		return
	}
	if self.cursorPoint.X == 1 {
		index := self.cursorPoint.Y - 1
		self.cursorPoint.X = len(self.text[index-1]) + 1
		self.text = append(
			self.text[:index-1],
			append(
				[][]Cell{append(self.text[index-1], self.text[index]...)},
				self.text[index+1:len(self.text)]...,
			)...,
		)
		self.cursorPoint.Y--
	} else {
		index := self.cursorPoint.Y - 1
		self.text[index] = append(
			self.text[index][:self.cursorPoint.X-2],
			self.text[index][self.cursorPoint.X-1:]...,
		)
		self.cursorPoint.X--
	}
}

// InsertText inserts the given text at the cursor position.
func (self *TextBox) InsertText(input string) {
	cells := ParseStyles(input, self.TextStyle)
	lines := SplitCells(cells, '\n')
	index := self.cursorPoint.Y - 1
	cellsAfterCursor := self.text[index][self.cursorPoint.X-1:]
	self.text[index] = append(self.text[index][:self.cursorPoint.X-1], lines[0]...)
	for i, line := range lines[1:] {
		index := self.cursorPoint.Y + i
		self.text = append(self.text[:index], append([][]Cell{line}, self.text[index:]...)...)
	}
	self.cursorPoint.Y += len(lines) - 1
	index = self.cursorPoint.Y - 1
	self.text[index] = append(self.text[index], cellsAfterCursor...)
	if len(lines) > 1 {
		self.cursorPoint.X = len(lines[len(lines)-1]) + 1
	} else {
		self.cursorPoint.X += len(lines[0])
	}
}

// ClearText clears the text and resets the cursor position.
func (self *TextBox) ClearText() {
	self.text = [][]Cell{[]Cell{}}
	self.cursorPoint = image.Pt(1, 1)
}

// SetText sets the text to the given text.
func (self *TextBox) SetText(input string) {
	self.ClearText()
	self.InsertText(input)
}

// GetText gets the text in string format along all its formatting tags
func (self *TextBox) GetText() string {
	return CellsToStyledString(JoinCells(self.text, '\n'), self.TextStyle)
}

// GetRawText gets the text in string format without any formatting tags
func (self *TextBox) GetRawText() string {
	return CellsToString(JoinCells(self.text, '\n'))
}

func (self *TextBox) MoveCursorLeft() {
	self.MoveCursor(self.cursorPoint.X-1, self.cursorPoint.Y)
}

func (self *TextBox) MoveCursorRight() {
	self.MoveCursor(self.cursorPoint.X+1, self.cursorPoint.Y)
}

func (self *TextBox) MoveCursorUp() {
	self.MoveCursor(self.cursorPoint.X, self.cursorPoint.Y-1)
}

func (self *TextBox) MoveCursorDown() {
	self.MoveCursor(self.cursorPoint.X, self.cursorPoint.Y+1)
}

func (self *TextBox) MoveCursor(x, y int) {
	self.cursorPoint.Y = MinInt(MaxInt(1, y), len(self.text))
	self.cursorPoint.X = MinInt(MaxInt(1, x), len(self.text[self.cursorPoint.Y-1])+1)
}

// CellsToStyledString converts []Cell to a string preserving the formatting tags
func CellsToStyledString(cells []Cell, defaultStyle Style) string {
	sb := strings.Builder{}
	runes := make([]rune, len(cells))
	currentStyle := cells[0].Style
	var j int

	for _, cell := range cells {
		if currentStyle != cell.Style {
			writeStyledText(&sb, runes[:j], currentStyle, defaultStyle)

			currentStyle = cell.Style
			j = 0
		}

		runes[j] = cell.Rune
		j++
	}

	// Write the last characters left in runes slice
	writeStyledText(&sb, runes[:j], currentStyle, defaultStyle)

	return sb.String()
}

func ContainsString(a []string, s string) bool {
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

//JoinCells converts [][]cell to a []cell using r as line breaker
func JoinCells(cells [][]Cell, r rune) []Cell {
	joinCells := make([]Cell, 0)
	lb := Cell{Rune: r, Style: StyleClear}
	length := len(cells)

	for i, cell := range cells {
		if i < length-1 {
			cell = append(cell, lb)
		}
		joinCells = append(joinCells, cell...)
	}

	return joinCells
}

func writeStyledText(sb *strings.Builder, runes []rune, currentStyle Style, defaultStyle Style) {
	if currentStyle != defaultStyle && currentStyle != StyleClear {
		sb.WriteByte('[')
		sb.WriteString(string(runes))
		sb.WriteByte(']')
		sb.WriteByte('(')
		sb.WriteString(StyleString(currentStyle))
		sb.WriteByte(')')
	} else {
		sb.WriteString(string(runes))
	}
}

//String returns a string representation of a Style
func StyleString(self Style) string {
	styles := make([]string, 0)

	if color, ok := textColorMap[self.Fg]; ok && self.Fg != StyleClear.Fg {
		styles = append(styles, tokenFg+tokenValueSeparator+color)
	}

	if color, ok := textColorMap[self.Bg]; ok && self.Bg != StyleClear.Bg {
		styles = append(styles, tokenBg+tokenValueSeparator+color)
	}

	if mod, ok := textModifierMap[self.Modifier]; ok && self.Modifier != StyleClear.Modifier {
		styles = append(styles, tokenModifier+tokenValueSeparator+mod)
	}

	return strings.Join(styles, tokenItemSeparator)
}
