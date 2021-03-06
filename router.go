package prelude

import "strings"

var (
	_logger *logger
)

func init() {
	_logger = &logger{
		mode: debugLevel,
	}
}

// HandlerFunc defines a function to server HTTP requests
type HandlerFunc func(c *Context) error

type Router struct {
	tree *tree
}

type tree struct {
	rootNode *node
}

type kind uint8

var (
	notFoundHandler = func(c *Context) {
		_logger.debug("NotFound")
	}
)

const (
	skind kind = iota
	pkind
	akind
)

// NewRouter function will create a new router instance
func NewRouter() *Router {
	return &Router{
		tree: &tree{
			rootNode: &node{
				parent:    nil,
				children:  []*node{},
				kind:      0,
				name:      "/",
				sortOrder: 0,
			},
		},
	}
}

// Invoke function is a middleware entry
// func (r *Router) Invoke(c *Context, next HandlerFunc) {
// 	h := r.Find(c.Request.URL.Path, c)

// 	if h == nil {
// 		next(c)
// 	} else {
// 		h(c)
// 	}
// }

// Add function which adding path and handler to router
func (r *Router) Add(path string, handler HandlerFunc) {
	_logger.debug("===Add")
	if len(path) == 0 {
		panic("router: path couldn't be empty")
	}

	_logger.debug("path:" + path)

	currentNode := r.tree.rootNode
	if path == "/" {
		currentNode.handler = handler
		return
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)
	pathParams := []string{}

	for index, element := range pathArray {
		if len(element) == 0 {
			continue
		}

		var childNode *node
		// this is static node
		childNode = currentNode.findChildByName(element)
		if childNode == nil {
			childNode = newNode(element, skind)
			currentNode.addChild(childNode)
			_logger.debug("add node")
		}

		// last node in the path
		if count == index+1 {
			childNode.params = pathParams
			childNode.handler = handler
		}

		currentNode = childNode
	}

}

// Find returns http handler for specific path
func (r *Router) Find(path string) HandlerFunc {
	_logger.debug("===Find")
	_logger.debug("path:" + path)
	if path[0] == '/' && len(path) > 1 {
		path = path[1:]
	}

	currentNode := r.tree.rootNode
	if path == "/" {
		return currentNode.handler
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)

	for index, element := range pathArray {
		// find static node first
		childNode := currentNode.findChildByName(element)

		if childNode == nil {
			//return notFoundHandler
			return nil
		}

		// last node in the path
		if count == index+1 {
			return childNode.handler
		}

		currentNode = childNode
	}

	//return notFoundHandler
	return nil
}

type node struct {
	parent    *node
	children  []*node
	kind      kind
	name      string
	pNames    []string
	params    []string
	sortOrder int
	handler   HandlerFunc
}

func newNode(name string, t kind) *node {
	return &node{
		kind:      t,
		name:      name,
		sortOrder: 0,
	}
}

func (n *node) addChild(node *node) {
	node.parent = n
	n.children = append(n.children, node)
}

func (n *node) findChildByName(name string) *node {
	var result *node
	for _, element := range n.children {
		if strings.EqualFold(element.name, name) && element.kind == skind {
			result = element
			break
		}
	}
	return result
}

func (n *node) findChildByKind(t kind) *node {
	for _, c := range n.children {
		if c.kind == t {
			return c
		}
	}
	return nil
}
