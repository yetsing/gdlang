package object

type Frame struct {
	filename string
	lineno   int
	funcName string
}

func (f *Frame) SetFilename(filename string) {
	f.filename = filename
}

func (f *Frame) GetFilename() string {
	return f.filename
}

func (f *Frame) SetLineno(lineno int) {
	f.lineno = lineno
}

func (f *Frame) GetLineno() int {
	return f.lineno
}

func (f *Frame) GetFuncName() string {
	return f.funcName
}

type CallStack struct {
	frames []*Frame
	index  int
}

func NewCallStack() *CallStack {
	return &CallStack{index: -1}
}

func (cs *CallStack) CreateFrame(filename string, funcName string) *Frame {
	frame := &Frame{filename: filename, funcName: funcName}
	cs.frames = append(cs.frames, frame)
	cs.index++
	return frame
}

func (cs *CallStack) DestroyFrame() *Frame {
	frame := cs.frames[cs.index]
	cs.frames = cs.frames[:cs.index]
	cs.index--
	return frame
}

func (cs *CallStack) Top() *Frame {
	return cs.frames[cs.index]
}

func (cs *CallStack) Copy() *CallStack {
	other := NewCallStack()
	other.frames = make([]*Frame, cs.index+1)
	other.index = cs.index
	copy(other.frames, cs.frames)
	return other
}

func (cs *CallStack) GetFrames() []*Frame {
	return cs.frames
}
