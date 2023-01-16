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
    local_globs int  // 下一个可用的局部变量插槽，从数组末尾开始
}

type Symble struct {
    Name string      // 标识符名
    Vartype Type     // 标识符类型

    IsLocal bool     // 是否是局部变量
    BelongFunc int   // 局部变量所属的函数
    Offset int       // 局部变量的偏移量

    EndLabel int     // 函数的末尾标签，用于return语句
    ReturnType Type  // 函数的返回类型
    FuncOffset int   // rsp栈顶的对齐偏移量
}

func init() {
    Gsym = &Symtable{
        symbles: make([]Symble, max_glob),
        globs:   0,
        local_globs: max_glob-1,
    }
}

// 查找全局符号name的插槽位置
func (s *Symtable) Findglob(name string) int {
    var i int
    for i = 0; i < s.globs; i++ {
        if s.symbles[i].Name == name {
            return i
        }
    }
    return -1
}

// 返回下一个可用的全局插槽位置
func (s *Symtable) Newglob() int {
    if s.globs >= s.local_globs {
        panic("Error: no available globs")
    }
    s.globs++
    return s.globs-1
}

// 新增一个全局符号到符号表
func (s *Symtable) Addglob(name string, vartype Type) int {
    var i int
    if i = s.Findglob(name); i != -1 {
        return i
    }

    i = s.Newglob()
    s.symbles[i].Name = name
    s.symbles[i].Vartype = vartype
    s.symbles[i].IsLocal = false
    s.symbles[i].FuncOffset = 0
    return i
}

func (s *Symtable) SetEndLabel(glob int, value int) {
    s.symbles[glob].EndLabel = value
}

func (s *Symtable) SetReturnType(glob int, value Type) {
    s.symbles[glob].ReturnType = value
}

func (s *Symtable) SetBelongFunc(glob int, value int) {
    s.symbles[glob].BelongFunc = value
}

func (s *Symtable) SetOffset(glob int, value int) {
    s.symbles[glob].Offset = value
}

func (s *Symtable) SetFuncOffset(glob int, value int) {
    s.symbles[glob].FuncOffset += value
}

func (s *Symtable) GetFuncOffset(glob int) int {
    offset := s.symbles[glob].FuncOffset
    mod := offset % 16
    if mod == 0 {
        return offset
    } else {
        return (offset / 16 + 1) * 16
    }
}

////////////////////////////////// 局部变量 ////////////////////////////
// 查找符号name的插槽位置
func (s *Symtable) Findlocal(name string) int {
    var i int
    for i = max_glob-1; i >= s.local_globs; i-- {
        if s.symbles[i].Name == name {
            return i
        }
    }
    return -1
}

// 返回下一个可用的插槽位置
func (s *Symtable) Newlocal() int {
    if s.globs >= s.local_globs {
        panic("Error: no available globs")
    }
    s.local_globs--
    return s.local_globs+1
}

// 新增一个符号到符号表
func (s *Symtable) Addlocal(name string, vartype Type) int {
    var i int
    if i = s.Findlocal(name); i != -1 {
        return i
    }

    i = s.Newlocal()
    s.symbles[i].Name = name
    s.symbles[i].Vartype = vartype
    s.symbles[i].IsLocal = true
    return i
}










