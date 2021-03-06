package main

import (
	"./excel"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

//filePath string, startStr string, endStr string, cellStart string, cellEnd string, separator string, colNum int, startRow int, sheetNum int
type MainWindows struct {
	*walk.MainWindow
	edit      *walk.TextEdit //输出文本框
	filePath  *walk.LineEdit //文件名输入框
	sheetNum  *walk.LineEdit //excel-sheet位置（开始为1，阿拉伯数字）
	colNum    *walk.LineEdit //excel-数据列列数（开始为1，阿拉伯数字）
	startRow  *walk.LineEdit //excel-起始行行数（开始为1，阿拉伯数字）
	startStr  *walk.LineEdit //拼接-开始文本
	cellStart *walk.LineEdit //拼接-单元格开始文本
	cellEnd   *walk.LineEdit //拼接-单元格结束文本
	separator *walk.LineEdit //拼接-分隔文本
	endStr    *walk.LineEdit //拼接-结束文本
}

func main() {
	mws := &MainWindows{}
	if err := (MainWindow{
		AssignTo: &mws.MainWindow,
		Title:    "Sql简单拼接",
		MinSize:  Size{600, 400},
		Size:     Size{1050, 750},
		MenuItems: []MenuItem{
			Menu{
				Text: "文件",
				Items: []MenuItem{
					Action{
						Text: "打开excel",
						Shortcut: Shortcut{ //定义快捷键后会有响应提示显示
							Modifiers: walk.ModControl,
							Key:       walk.KeyO,
						},
						OnTriggered: mws.openFileActionTriggered, //点击动作触发响应函数
					},
					Action{
						Text: "导出sql",
						Shortcut: Shortcut{
							Modifiers: walk.ModControl | walk.ModShift,
							Key:       walk.KeyS,
						},
						OnTriggered: mws.saveFileActionTriggered,
					},
					Action{
						Text: "退出",
						OnTriggered: func() {
							_ = mws.Close()
						},
					},
				},
			},
			Menu{
				Text: "帮助",
				Items: []MenuItem{
					Action{
						Text: "关于",
						OnTriggered: func() {
							walk.MsgBox(mws, "关于", "这个工具为王sao峰设计，疫情项目每天一点两点要转sql。另外《武汉战疫小程序》UP！！！--武汉战疫项目组-0x5f81",
								walk.MsgBoxIconInformation|walk.MsgBoxDefButton1)
						},
					},
				},
			},
		},

		//sheetNum  *walk.LineEdit //excel-sheet位置（开始为1，阿拉伯数字）
		//colNum    *walk.LineEdit //excel-数据列列数（开始为1，阿拉伯数字）
		//startRow  *walk.LineEdit //excel-起始行行数（开始为1，阿拉伯数字）

		Layout: VBox{},
		Children: []Widget{
			Label{Text: "excel路径"},
			LineEdit{AssignTo: &mws.filePath},
			Label{Text: "excel-sheet位置（开始为1，阿拉伯数字）"},
			LineEdit{AssignTo: &mws.sheetNum, Text: "1"},
			Label{Text: "excel-数据列列数（开始为1，阿拉伯数字）"},
			LineEdit{AssignTo: &mws.colNum, Text: "7"},
			Label{Text: "excel-起始行行数（开始为1，阿拉伯数字）"},
			LineEdit{AssignTo: &mws.startRow, Text: "1"},
			//startStr  *walk.LineEdit //拼接-开始文本
			Label{Text: "拼接-开始文本"},
			LineEdit{AssignTo: &mws.startStr, Text: "in ("},
			//cellStart *walk.LineEdit //拼接-单元格开始文本
			Label{Text: "拼接-单元格开始文本"},
			LineEdit{AssignTo: &mws.cellStart, Text: "'"},
			//cellEnd   *walk.LineEdit //拼接-单元格结束文本
			Label{Text: "拼接-单元格结束文本"},
			LineEdit{AssignTo: &mws.cellEnd, Text: "'"},
			//separator *walk.LineEdit //拼接-分隔文本
			Label{Text: "拼接-分隔文本"},
			LineEdit{AssignTo: &mws.separator, Text: ","},
			//endStr    *walk.LineEdit //拼接-结束文本
			Label{Text: "拼接-结束文本"},
			LineEdit{AssignTo: &mws.endStr, Text: ")"},
			PushButton{
				Column:     1,
				ColumnSpan: 1,
				Text:       "生产拼接文本",
				OnClicked:  mws.ProcessXlsxToSq,
			},
			TextEdit{
				AssignTo: &mws.edit,
			},
		},
		OnDropFiles: mws.dropFiles,
	}).Create(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	mws.Run()
}

func (mws *MainWindows) Alert(title string, content string) {
	walk.MsgBox(mws, title, content, walk.MsgBoxIconWarning)
}

func (mws *MainWindows) ProcessXlsxToSq() {
	var err error
	var sqlStr string
	sheetNum := mws.sheetNum.Text()
	sheetNumInt, err := strconv.Atoi(sheetNum)
	if err != nil {
		mws.Alert("警告", "excel-sheet位置必须是开始为1的阿拉伯数字")
		return
	}
	colNum := mws.colNum.Text()
	colNumInt, err := strconv.Atoi(colNum)
	if err != nil {
		mws.Alert("警告", "excel-数据列列数必须是开始为1的阿拉伯数字")
		return
	}
	startRow := mws.startRow.Text()
	startRowInt, err := strconv.Atoi(startRow)
	if err != nil {
		mws.Alert("警告", "excel-起始行行数必须是开始为1的阿拉伯数字")
		return
	}
	var ext excel.ExTools
	if path.Ext(mws.filePath.Text()) == ".xls" {
		ext = excel.XlsTools{
			FilePath:  mws.filePath.Text(),
			StartStr:  mws.startStr.Text(),
			EndStr:    mws.endStr.Text(),
			CellStart: mws.cellStart.Text(),
			CellEnd:   mws.cellEnd.Text(),
			Separator: mws.separator.Text(),
			ColNum:    colNumInt,
			StartRow:  startRowInt,
			SheetNum:  sheetNumInt,
		}
	} else if path.Ext(mws.filePath.Text()) == ".xlsx" {
		ext = excel.XlsxTools{
			FilePath:  mws.filePath.Text(),
			StartStr:  mws.startStr.Text(),
			EndStr:    mws.endStr.Text(),
			CellStart: mws.cellStart.Text(),
			CellEnd:   mws.cellEnd.Text(),
			Separator: mws.separator.Text(),
			ColNum:    colNumInt,
			StartRow:  startRowInt,
			SheetNum:  sheetNumInt,
		}
	} else {
		mws.Alert("警告", "解析文件后缀名需要是xls或xlsx，目前不支持其他格式。（可以用过其他工具另存为以上两种格式）|"+path.Ext(mws.filePath.Text()))
		return
	}

	if sqlStr, err = ext.ParseToSql(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(sqlStr)
	}
	if err != nil {
		mws.Alert("警告", err.Error())
		return
	}
	_ = mws.edit.SetText(sqlStr)
}

//打开文件方案
func (mws *MainWindows) openFileActionTriggered() {
	dlg := new(walk.FileDialog)
	dlg.Title = "打开xlsx"
	dlg.Filter = "excel (*.xlsx)|*.xlsx|excel (*.xls)|*.xls"

	if ok, err := dlg.ShowOpen(mws); err != nil {
		fmt.Fprintln(os.Stderr, "错误：打开文件时\r\n")
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "用户取消\r\n")
		return
	}
	mws.filePath.SetText(dlg.FilePath)
}

func (mws *MainWindows) saveFileActionTriggered() {
	dlg := new(walk.FileDialog)
	dlg.Title = "导出"

	if ok, err := dlg.ShowSave(mws); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "取消")
		return
	}

	data := mws.edit.Text()
	filename := dlg.FilePath
	f, err := os.Open(filename)
	if err != nil {
		f, _ = os.Create(filename)
	} else {
		f.Close()
		f, err = os.OpenFile(filename, os.O_WRONLY, 0x666)
	}
	if len(data) == 0 {
		f.Close()
		return
	}
	io.Copy(f, strings.NewReader(data))
	f.Close()
}

func (mws *MainWindows) newAction_Triggered() {
	walk.MsgBox(mws, "New", "Newing something up... or not.", walk.MsgBoxIconInformation)
}

func (mws *MainWindows) changeViewAction_Triggered() {
	walk.MsgBox(mws, "Change View", "By now you may have guessed it. Nothing changed.", walk.MsgBoxIconInformation)
}

func (mws *MainWindows) dropFiles(files []string) {
	mws.edit.SetText("")
	for _, v := range files {
		mws.edit.AppendText(v + "\r\n")
	}
}
