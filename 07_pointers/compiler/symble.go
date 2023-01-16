package compiler

var Gsym *Symtable
const max_glob = 1024

type Type int
const (
    VAR_CHAR Type = iota
    VAR_INT
    VAR_FLOAT
    VAR_STRING
    VAR_POINTER_CHAR
    VAR_POINTER_INT
    VAR_ARRAY
    VAR_STRCUT
    VAR_INTERFACE
    VAR_FUNC
)

type Symtable struct {
    symbles []Symble
    globs   int  // 下一个可用的插槽
}

type Symble struct {
    Name string      // 标识符名
    Vartype Type     // 标识符类型

    EndLabel int     // 函数的末尾标签，用于return语句
    ReturnType Type  // 函数的返回类型
}

func init() {
    Gsym = &Symtable{
        symbles: make([]Symble, max_glob),
        globs:   0,
    }
}

// 查找符号name的插槽位置
func (s *Symtable) Findglob(name string) int {
    var i int
    for i = 0; i < s.globs; i++ {
        if s.symbles[i].Name == name {
            return i
        }
    }
    return -1
}

// 返回下一个可用的插槽位置
func (s *Symtable) Newglob() int {
    if s.globs >= max_glob {
        panic("Error: no available globs")
    }
    s.globs++
    return s.globs-1
}

// 新增一个符号到符号表
func (s *Symtable) Addglob(name string, vartype Type) int {
    var i int
    if i = s.Findglob(name); i != -1 {
        return i
    }

    i = s.Newglob()
    s.symbles[i].Name = name
    s.symbles[i].Vartype = vartype
    return i
}

func (s *Symtable) Setglob(glob int, value interface{}) {
    switch value.(type) {
    case int:
        s.symbles[glob].EndLabel = value.(int)
    case Type:
        s.symbles[glob].ReturnType = value.(Type)
    }
}








