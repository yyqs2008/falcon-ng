package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/errors"
	"github.com/toolkits/pkg/str"

	"github.com/open-falcon/falcon-ng/src/model"
)

type nodeForm struct {
	Pid  int64  `json:"pid"`
	Name string `json:"name"`
	Leaf int    `json:"leaf"`
	Note string `json:"note"`
}

func nodePost(c *gin.Context) {
	var f nodeForm
	errors.Dangerous(c.ShouldBind(&f))

	if str.Dangerous(f.Name) {
		errors.Bomb("name invalid")
	}

	if str.Dangerous(f.Note) {
		errors.Bomb("note invalid")
	}

	if f.Pid <= 0 {
		errors.Bomb("pid invalid")
	}

	if f.Leaf != 0 && f.Leaf != 1 {
		errors.Bomb("leaf invalid")
	}

	if !str.IsMatch(f.Name, `^[a-zA-Z0-9\-]+$`) {
		errors.Bomb("name valid characters: [a-zA-Z0-9] and -")
	}

	parent, err := model.NodeGet("id", f.Pid)
	errors.Dangerous(err)

	if parent == nil {
		errors.Bomb("arg[pid] invalid, no such parent node")
	}

	if parent.Type > 0 {
		errors.Bomb("由其他子系统管理的节点，不允许在服务树视图添加子节点")
	}

	newPath := parent.Path + "." + f.Name
	node, err := model.NodeGet("path", newPath)
	errors.Dangerous(err)

	if node != nil {
		errors.Bomb("%s already exists", newPath)
	}

	_, err = parent.CreateChild(f.Name, f.Leaf, 0, f.Note)
	renderMessage(c, err)
}

func nodeSearchGet(c *gin.Context) {
	limit := queryInt(c, "limit", 20)
	query := queryStr(c, "query", "")
	nodes, err := model.NodeQueryPath(query, limit)
	renderData(c, nodes, err)
}

type nodeNameForm struct {
	Name string `json:"name" binding:"required"`
}

func nodeNamePut(c *gin.Context) {
	var f nodeNameForm
	errors.Dangerous(c.ShouldBind(&f))

	if !str.IsMatch(f.Name, `^[a-zA-Z0-9\-]+$`) {
		errors.Bomb("name valid characters: [a-zA-Z0-9] and -")
	}

	node, err := model.NodeGet("id", urlParamInt64(c, "id"))
	errors.Dangerous(err)

	if node == nil {
		errors.Bomb("arg[id] invalid, no such node")
	}

	if node.Type > 0 {
		errors.Bomb("由其他子系统管理的节点，不允许在服务树视图改名")
	}

	renderMessage(c, node.Rename(f.Name))
}

func nodeDel(c *gin.Context) {
	node, err := model.NodeGet("id", urlParamInt64(c, "id"))
	errors.Dangerous(err)

	if node == nil {
		errors.Bomb("arg[id] invalid, no such node")
	}

	if node.Type > 0 {
		errors.Bomb("由其他子系统管理的节点，不允许在服务树视图删除")
	}

	renderMessage(c, node.Del())
}

func nodeLeafIdsGet(c *gin.Context) {
	idsStr := mustQueryStr(c, "ids")
	ids := str.IdsInt64(idsStr)

	nodes, err := model.NodesGetByIds(ids)
	errors.Dangerous(err)

	dat := make(map[int64][]int64)

	for i := 0; i < len(nodes); i++ {
		leafIds, err := nodes[i].LeafIds()
		errors.Dangerous(err)
		dat[nodes[i].Id] = leafIds
	}

	renderData(c, dat, nil)
}

func nodePidsGet(c *gin.Context) {
	idsStr := mustQueryStr(c, "ids")
	ids := str.IdsInt64(idsStr)

	nodes, err := model.NodesGetByIds(ids)
	errors.Dangerous(err)

	dat := make(map[int64][]int64)

	for i := 0; i < len(nodes); i++ {
		pids, err := nodes[i].Pids()
		errors.Dangerous(err)
		dat[nodes[i].Id] = pids
	}

	renderData(c, dat, err)
}

func nodesByIdsGets(c *gin.Context) {
	idsStr := mustQueryStr(c, "ids")
	ids := str.IdsInt64(idsStr)
	nodes, err := model.NodeByIds(ids)
	renderData(c, nodes, err)
}
