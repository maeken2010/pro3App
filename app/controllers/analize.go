package controllers

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"
  // "github.com/revel/revel"
  // "pro3App/app/chatroom"
  // "pro3App/app/models"
)

// func main() {
//   lines := fromFile("test_room.txt")
//   // fmt.Printf("lines: %v\n", analize("0+0+1",lines))
//   request := os.Args[1]
//   fmt.Println(analize(request,lines))
// }

// ex)req="0+1+1"
func analize(req string,path string) string{
  roomData := fromFile(path)
  if len(req)==0 {
    fmt.Fprintf(os.Stderr, "request error!")
  }
  reqList := strings.FieldsFunc(req,isPlusRune) //reqを'+'で区切ってlistに変換
  m := 0
  for i := 0; i < len(reqList); i++ {
    // levelLines := make([]string, 0)
    x := 0
    for  {
      if isHeadSpace(roomData[m],i) {
        if strconv.Itoa(x)==reqList[i] {
          break
        }
        x++
      }
      m++
    }
  }
  return strings.Trim(roomData[m]," ")
  // return reqList
}

func isPlusRune(s rune) bool {
    return s=='+'
}

func isHeadSpace(s string,n int) bool{
  if n==0 {
    return s[n] != ' '
  }
  return s[n-1]==' ' && s[n]!=' '
}


func fromFile(filePath string) []string {
  // ファイルを開く
  f, err := os.Open(filePath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filePath, err)
    os.Exit(1)
  }

  // 関数return時に閉じる
  defer f.Close()

  // Scannerで読み込む
  // lines := []string{}
  lines := make([]string, 0, 100)  // ある程度行数が事前に見積もれるようであれば、makeで初期capacityを指定して予めメモリを確保しておくことが望ましい
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    // appendで追加
    lines = append(lines, scanner.Text())
  }
  if serr := scanner.Err(); serr != nil {
    fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
  }

  return lines
}
