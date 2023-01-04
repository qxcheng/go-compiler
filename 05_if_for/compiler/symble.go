package compiler

var Gsym *symtable
const max_glob = 1024

type symtable struct {
    symbles []symble
    globs   int  // 下一个可用的插槽
}

type symble struct {
    name string
}

func init() {
    Gsym = &symtable{
        symbles: make([]symble, max_glob),
        globs:   0,
    }
}

// 查找符号name的插槽位置
func (s *symtable) Findglob(name string) int {
    var i int
    for i = 0; i < s.globs; i++ {
        if s.symbles[i].name == name {
            return i
        }
    }
    return -1
}

// 返回下一个可用的插槽位置
func (s *symtable) Newglob() int {
    if s.globs >= max_glob {
        panic("Error: no available globs")
    }
    s.globs++
    return s.globs-1
}

// 新增一个符号到符号表
func (s *symtable) Addglob(name string) int {
    var i int
    if i = s.Findglob(name); i != -1 {
        return i
    }

    i = s.Newglob()
    s.symbles[i].name = name
    return i
}








