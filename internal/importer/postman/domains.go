package postman

var PostManCollection string

type PostmanCollection struct {
	Item []Item `json:"item"`
}

type Item struct {
	Name    string  `json:"name"`
	Item    []Item  `json:"item"`
	Request Request `json:"request"`
}

type Request struct {
	Method string `json:"method"`
	Url    Url    `json:"url"`
}

type Url struct {
	Raw      string     `json:"raw"`
	Variable []Variable `json:"variable"`
}

type Variable struct {
	Key string `json:"key"`
}

type Handler struct {
	name         string
	method       string
	fields       any
	pathVariable any
}

func (h *Handler) GetName() string {
	return h.name
}

func (h *Handler) SetName(name string) {
	h.name = name
}

func (h *Handler) GetMethod() string {
	return h.method
}

func (h *Handler) SetMethod(method string) {
	h.method = method
}

func (h *Handler) GetFields() any {
	return h.fields
}

func (h *Handler) SetFields(fields any) {
	h.fields = fields
}

func (h *Handler) GetPathVariable() any {
	return h.pathVariable
}

func (h *Handler) SetPathVariable(queries any) {
	h.pathVariable = queries
}
